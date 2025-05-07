using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Authorization;
using System.Text.Json;
using System.Data;
using Dapper;
using perspectize_be.DTOs;
using perspectize_be.Models;
using perspectize_be.Services;
using perspectize_be.DTOs.perspectize_be.DTOs;

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
        
        [HttpPost("video")]
        public IActionResult PostVideo([FromBody] VideoRequest req, CancellationToken cancellationToken)
        {
            if (string.IsNullOrEmpty(req.VideoId)) return BadRequest(new { message = "videoId is required" });
            return Ok(
                new { message = $"Video '{req.VideoId}' posted successfully" }
            );
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
                    JsonElement videoData = JsonSerializer.Deserialize<JsonElement>(responseContent);
                    
                    if (!videoData.TryGetProperty("items", out JsonElement items) || items.GetArrayLength() == 0)
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
                    
                    // Using Dapper to query existing content
                    string findQuery = "SELECT * FROM content WHERE url = @Url LIMIT 1";
                    Content? existingContent = await _dbConnection.QueryFirstOrDefaultAsync<Content>(findQuery, new { Url = url });
                    
                    if (existingContent != null)
                    {
                        // Using Dapper for update
                        string updateQuery = @"
                            UPDATE content 
                            SET length = @Length, 
                                length_units = @LengthUnits, 
                                response = @Response::jsonb, 
                                name = @Name, 
                                updated_at = @UpdatedAt 
                            WHERE url = @Url";
                        
                        await _dbConnection.ExecuteAsync(updateQuery, new { 
                            Length = durationInSeconds.ToString(),
                            LengthUnits = "seconds",
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
                        // Using Dapper for insert
                        string insertQuery = @"
                            INSERT INTO content (url, length, length_units, response, content_type, name, created_at, updated_at)
                            VALUES (@Url, @Length, @LengthUnits, @Response::jsonb, @ContentType, @Name, @CreatedAt, @UpdatedAt)";
                        
                        await _dbConnection.ExecuteAsync(insertQuery, new {
                            Url = url,
                            Length = durationInSeconds.ToString(),
                            LengthUnits = "seconds",
                            Response = responseContent,
                            ContentType = "youtube",
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

        [HttpPut("videos")]
        public async Task<IActionResult> PutVideos([FromBody] VideosRequest request, CancellationToken cancellationToken)
        {
            return await PostVideos(request, cancellationToken);
        }
    }
}