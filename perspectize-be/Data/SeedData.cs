using Microsoft.EntityFrameworkCore;
using perspectize_be.Data;
using perspectize_be.Models;
using perspectize_be.Services;
using System.Text.Json;

namespace perspectize_be.Data
{
    public static class SeedData
    {
        public static async Task InitializeAsync(IServiceProvider serviceProvider)
        {
            using var scope = serviceProvider.CreateScope();
            ApplicationDbContext context = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
            IConfiguration configuration = scope.ServiceProvider.GetRequiredService<IConfiguration>();
            IHttpClientFactory httpClientFactory = scope.ServiceProvider.GetRequiredService<IHttpClientFactory>();

            // Create YouTube service
            HttpClient httpClient = httpClientFactory.CreateClient();
            YouTubeService youtubeService = new YouTubeService(httpClient, configuration);

            // Check if we already have content
            if (await context.Contents.AnyAsync())
            {
                return; // DB has already been seeded
            }

            // List of interesting YouTube videos to seed
            List<string> videoUrls = new List<string>
            {
                "https://www.youtube.com/watch?v=dQw4w9WgXcQ", // Rick Astley - Never Gonna Give You Up
                "https://www.youtube.com/watch?v=9bZkp7q19f0", // PSY - Gangnam Style
                "https://www.youtube.com/watch?v=PGNiXGX2nLU", // Dead or Alive - You Spin Me Round
                "https://www.youtube.com/watch?v=W6DmHGYy_xk", // "Evolution of Dance" - Funny dance compilation
                "https://www.youtube.com/watch?v=MtN1YnoL46Q", // "The Duck Song"
                "https://www.youtube.com/watch?v=ApLAHmcaSZM", // "Dude Perfect: Water Bottle Flip Edition"
                "https://www.youtube.com/watch?v=qxVDOpOEFjI", // "Ultimate Dog Tease"
                "https://www.youtube.com/watch?v=ZyhrYis509A", // Aqua - Barbie Girl
                "https://www.youtube.com/watch?v=KQ6zr6kCPj8", // LMFAO - Party Rock Anthem
                "https://www.youtube.com/watch?v=QH2-TGUlwu4", // Nyan Cat
                "https://www.youtube.com/watch?v=lsJLLEwUYZM", // "Kid Snippets: Math Class"
                "https://www.youtube.com/watch?v=4WX58CZwyiU", // "Useless Machine"
                "https://www.youtube.com/watch?v=z1Kdoja3hlk", // JavaScript: Understanding the Weird Parts
                "https://www.youtube.com/watch?v=X9eRLElSW1c", // C# Advanced Tutorial
                "https://www.youtube.com/watch?v=8jLOx1hD3_o", // ASP.NET Core Crash Course
                "https://www.youtube.com/watch?v=nx2-4l4s4Nw", // Dependency Injection in .NET
                "https://www.youtube.com/watch?v=QQVAHbKdcKM", // YouTube API in C#
                "https://www.youtube.com/watch?v=IdG-rF72z4s", // Building REST APIs with .NET
                "https://www.youtube.com/watch?v=wgXwRya1BNk", // JWT Authentication in .NET
                "https://www.youtube.com/watch?v=JeVYfpQvxPY"  // React with ASP.NET Core
            };

            foreach (string url in videoUrls)
            {
                try
                {
                    // Extract video ID from URL
                    string videoId = youtubeService.ExtractVideoId(url);

                    // Get video details from YouTube API
                    string getUrl = $"https://www.googleapis.com/youtube/v3/videos?key={youtubeService.GetApiKey()}&part=snippet,contentDetails&id={videoId}";
                    HttpResponseMessage response = await httpClient.GetAsync(getUrl);

                    if (!response.IsSuccessStatusCode)
                    {
                        Console.WriteLine($"Error getting video {url}: YouTube API returned status code {response.StatusCode}");
                        continue;
                    }

                    string responseContent = await response.Content.ReadAsStringAsync();
                    JsonElement videoData = JsonSerializer.Deserialize<JsonElement>(responseContent);

                    // Check if the video exists
                    if (!videoData.TryGetProperty("items", out JsonElement items) || items.GetArrayLength() == 0)
                    {
                        Console.WriteLine($"Error getting video {url}: Video not found");
                        continue;
                    }

                    JsonElement videoItem = items[0];

                    // Get video details for content entry
                    string? title = videoItem.GetProperty("snippet").GetProperty("title").GetString();
                    string? durationISO = videoItem.GetProperty("contentDetails").GetProperty("duration").GetString();
                    int? durationInSeconds = youtubeService.ConvertDurationToSeconds(durationISO ?? string.Empty);

                    // Create new content
                    Content newContent = new Content
                    {
                        Url = url,
                        Length = durationInSeconds.ToString(),
                        LengthUnits = "seconds",
                        Response = JsonDocument.Parse(responseContent),
                        ContentType = "youtube",
                        Name = title ?? string.Empty,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow
                    };

                    context.Contents.Add(newContent);
                    Console.WriteLine($"Added video: {title}");
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error seeding video {url}: {ex.Message}");
                }
            }

            await context.SaveChangesAsync();
            Console.WriteLine("Seed data completed");
        }
    }
}