// Controllers/YTController.cs
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using System.Threading;
using Microsoft.AspNetCore.Authorization;
using System.Text.Json;
using perspectize_be.Data;
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
        private readonly ApplicationDbContext _context;
        private readonly YouTubeService _youtubeService;

        public YTController(ApplicationDbContext context, YouTubeService youtubeService)
        {
            _context = context;
            _youtubeService = youtubeService;
        }

        // Endpoint 1: GET video - implemented as in postman API collection
        [HttpGet("video")]
        public async Task<IActionResult> GetVideo([FromQuery] string videoId, CancellationToken cancellationToken)
        {
            if (string.IsNullOrEmpty(videoId))
            {
                return BadRequest(new { message = "videoId is required" });
            }

            try
            {
                // Following the Postman collection format
                string url = $"https://www.googleapis.com/youtube/v3/videos?key={_youtubeService.GetApiKey()}&part=snippet,contentDetails&id={videoId}";
                HttpResponseMessage response = await _youtubeService.HttpClient.GetAsync(url, cancellationToken);
                
                if (!response.IsSuccessStatusCode)
                {
                    string errorContent = await response.Content.ReadAsStringAsync(cancellationToken);
                    return StatusCode((int)response.StatusCode, new { message = $"YouTube API returned status code {response.StatusCode}: {errorContent}" });
                }
                
                string content = await response.Content.ReadAsStringAsync(cancellationToken);
                
                // Return the response directly from YouTube API
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

        // Endpoint 2: POST videos - takes an array of YouTube URLs
        [HttpPost("videos")]
        public async Task<IActionResult> PostVideos([FromBody] VideosRequest request, CancellationToken cancellationToken)
        {
            if (request?.VideoUrls == null || !request.VideoUrls.Any())
            {
                return BadRequest(new { message = "At least one video URL is required" });
            }

            List<object> results = new List<object>();

            foreach (var url in request.VideoUrls)
            {
                try
                {
                    // Extract video ID from URL
                    string videoId = _youtubeService.ExtractVideoId(url);
                    
                    // Intermediate call to the GET endpoint (YouTube API)
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
                    
                    // Check if the video exists
                    if (!videoData.TryGetProperty("items", out var items) || items.GetArrayLength() == 0)
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
                    
                    // Get video details for content entry
                    string? title = videoItem.GetProperty("snippet").GetProperty("title").GetString();
                    string? durationISO = videoItem.GetProperty("contentDetails").GetProperty("duration").GetString();
                    int? durationInSeconds = _youtubeService.ConvertDurationToSeconds(durationISO ?? string.Empty);
                    
                    // Check if content already exists with this URL - implement upsert
                    Content? existingContent = await _context.Contents
                        .FirstOrDefaultAsync(c => c.Url == url, cancellationToken);
                    
                    if (existingContent != null)
                    {
                        // Update existing content
                        existingContent.Length = durationInSeconds.ToString();
                        existingContent.LengthUnits = "seconds";
                        existingContent.Response = videoData;
                        existingContent.Name = title ?? string.Empty;
                        existingContent.UpdatedAt = DateTime.UtcNow;
                        
                        _context.Contents.Update(existingContent);
                        
                        results.Add(new
                        {
                            status = "updated",
                            videoId,
                            name = existingContent.Name,
                            url = existingContent.Url
                        });
                    }
                    else
                    {
                        // Create new content
                        Content newContent = new Content
                        {
                            Url = url,
                            Length = durationInSeconds.ToString(),
                            LengthUnits = "seconds",
                            Response = videoData,
                            ContentType = "youtube",
                            Name = title ?? string.Empty,
                            CreatedAt = DateTime.UtcNow,
                            UpdatedAt = DateTime.UtcNow
                        };
                        
                        _context.Contents.Add(newContent);
                        
                        results.Add(new
                        {
                            status = "created",
                            videoId,
                            name = newContent.Name,
                            url = newContent.Url
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

            // Save all changes to the database
            await _context.SaveChangesAsync(cancellationToken);
            
            return Ok(results);
        }

        // Endpoint 3: PUT videos - same as POST
        [HttpPut("videos")]
        public async Task<IActionResult> PutVideos([FromBody] VideosRequest request, CancellationToken cancellationToken)
        {
            // Same implementation as POST for upsert functionality
            return await PostVideos(request, cancellationToken);
        }
    }
}