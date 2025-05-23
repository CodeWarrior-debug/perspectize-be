using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Authorization;
using System.Text.Json;
using System.Data;
using Dapper;
using perspectize_be.DTOs;
using perspectize_be.Models;
using perspectize_be.Services;
using perspectize_be.Common;
namespace perspectize_be.Controllers
{
    [Route("youtube")]
    [ApiController]
    [AllowAnonymous]
    public class YTController : ControllerBase
    {
        private readonly YouTubeService _youtubeService;
        private readonly IDbConnection _dbConnection;

        public YTController(YouTubeService youtubeService, IDbConnection dbConnection)
        {
            _youtubeService = youtubeService;
            _dbConnection = dbConnection;
        }

        [HttpGet("video")]
        public async Task<IActionResult> GetVideo([FromQuery] string videoId, CancellationToken cancellationToken)
        {
            if (string.IsNullOrEmpty(videoId))
            {
                return BadRequest(new { message = "videoId is required" });
            }

            try
            {
                string url = $"https://www.googleapis.com/youtube/v3/videos?key={_youtubeService.GetApiKey()}&part=snippet,contentDetails&id={videoId}";
                HttpResponseMessage response = await _youtubeService.HttpClient.GetAsync(url, cancellationToken);
                
                if (!response.IsSuccessStatusCode)
                {
                    string errorContent = await response.Content.ReadAsStringAsync(cancellationToken);
                    return StatusCode((int)response.StatusCode, new { message = $"YouTube API returned status code {response.StatusCode}: {errorContent}" });
                }
                
                string content = await response.Content.ReadAsStringAsync(cancellationToken);
                
                return Content(content, "application/json");
            }
            catch (Exception ex)
            {
                return StatusCode(500, new { message = $"Error retrieving video: {ex.Message}" });
            }
        }
        
        [HttpPost("videos")]
        public async Task<IActionResult> PostVideos([FromBody] VideosRequest request, CancellationToken cancellationToken)
        {
            if (request?.VideoUrls == null || !request.VideoUrls.Any())
            {
                return BadRequest(new { message = "At least one video URL is required" });
            }

            List<object> results = new List<object>();

            foreach (string url in request.VideoUrls)
            {
                try
                {
                    string videoId = _youtubeService.ExtractVideoId(url);
                    
                    string getUrl = $"https://www.googleapis.com/youtube/v3/videos?key={_youtubeService.GetApiKey()}&part=snippet,contentDetails&id={videoId}";
                    HttpResponseMessage response = await _youtubeService.HttpClient.GetAsync(getUrl, cancellationToken);
                    
                    if (!response.IsSuccessStatusCode)
                    {
                        results.Add(new
                        {
                            status = "error",
                            url,
                            message = $"YouTube API returned status code {response.StatusCode}"
                        });
                        continue;
                    }
                    
                    string responseContent = await response.Content.ReadAsStringAsync(cancellationToken);
                    
                    if (!responseContent.StartsWith("{"))
                    {
                        results.Add(new
                        {
                            status = "error",
                            url,
                            message = "Invalid response format"
                        });
                        continue;
                    }
                    
                    JsonDocument videoData = JsonDocument.Parse(responseContent);
                    
                    if (!videoData.RootElement.TryGetProperty("items", out JsonElement items) || items.GetArrayLength() == 0)
                    {
                        results.Add(new
                        {
                            status = "error",
                            url,
                            message = "Video not found"
                        });
                        continue;
                    }
                    
                    JsonElement videoItem = items[0];
                    
                    string? title = videoItem.GetProperty("snippet").GetProperty("title").GetString();
                    string? durationISO = videoItem.GetProperty("contentDetails").GetProperty("duration").GetString();
                    int? durationInSeconds = _youtubeService.ConvertDurationToSeconds(durationISO ?? string.Empty);
                    
                    string findQuery = "SELECT * FROM content WHERE url = @Url LIMIT 1"; //TODO: refactor to ON CONFLICT upsert method (keep existing for now)
                    Content? existingContent = await _dbConnection.QueryFirstOrDefaultAsync<Content>(findQuery, new { Url = url });
                    
                    if (existingContent != null)
                    {
                        string updateQuery = @"
                            UPDATE content 
                            SET length = @Length, 
                                length_units = @LengthUnits, 
                                response = @Response::jsonb, 
                                name = @Name, 
                                updated_at = @UpdatedAt 
                            WHERE url = @Url";
                        
                        await _dbConnection.ExecuteAsync(updateQuery, new { 
                            Length = durationInSeconds,
                            LengthUnits = Constants.LengthUnits.Seconds,
                            Response = responseContent,
                            Name = title ?? string.Empty,
                            UpdatedAt = DateTime.UtcNow,
                            Url = url
                        });
                        
                        results.Add(new
                        {
                            status = "updated",
                            videoId,
                            name = title ?? string.Empty,
                            url = url
                        });
                    }
                    else
                    {
                        string insertQuery = @"
                            INSERT INTO content (url, length, length_units, response, content_type, name, created_at, updated_at)
                            VALUES (@Url, @Length, @LengthUnits, @Response::jsonb, @ContentType, @Name, @CreatedAt, @UpdatedAt)";
                        
                        await _dbConnection.ExecuteAsync(insertQuery, new {
                            Url = url,
                            Length = durationInSeconds,
                            LengthUnits = Constants.LengthUnits.Seconds,
                            Response = responseContent,
                            ContentType = Constants.ContentType.YouTube,
                            Name = title ?? string.Empty,
                            CreatedAt = DateTime.UtcNow,
                            UpdatedAt = DateTime.UtcNow
                        });
                        
                        results.Add(new
                        {
                            status = "created",
                            videoId,
                            name = title ?? string.Empty,
                            url = url
                        });
                    }
                }
                catch (Exception ex)
                {
                    results.Add(new
                    {
                        status = "error",
                        url,
                        message = ex.Message
                    });
                }
            }
            
            return Ok(results);
        }

        [HttpPut("videos")] //TODO: refactor later - this is clever reuse, but uses the more expensive GET then INSERT / UPDATE approach. What we're looking for is just the simple Update on one video here.
        public async Task<IActionResult> PutVideos([FromBody] VideosRequest request, CancellationToken cancellationToken)
        {
            return await PostVideos(request, cancellationToken);
        }
    }
}