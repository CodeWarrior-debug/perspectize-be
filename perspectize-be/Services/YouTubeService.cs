// Services/YouTubeService.cs
using System.Text.Json;
using System.Text.RegularExpressions;
using System.Xml;

namespace perspectize_be.Services
{
    public class YouTubeService
    {
        private readonly IConfiguration _configuration;
        private readonly string _apiKey;
        
        public HttpClient HttpClient { get; }
        
        public string GetApiKey() => _apiKey;

        public YouTubeService(HttpClient httpClient, IConfiguration configuration)
        {
            HttpClient = httpClient;
            _configuration = configuration;
            _apiKey = _configuration["YouTube:ApiKey"] ?? throw new ArgumentNullException("YouTube:ApiKey", "YouTube API key is missing in configuration");
        }

        // Extract video ID from YouTube URL
        public string ExtractVideoId(string url)
        {
            // Handle different YouTube URL formats
            Regex youtubeIdRegex = new Regex(@"(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^""&?\/\s]{11})");
            Match match = youtubeIdRegex.Match(url);
            
            if (match.Success)
            {
                return match.Groups[1].Value;
            }
            
            throw new ArgumentException($"Could not extract YouTube video ID from URL: {url}");
        }

        // Convert ISO 8601 duration to seconds
        public int ConvertDurationToSeconds(string isoDuration)
        {
            TimeSpan duration = XmlConvert.ToTimeSpan(isoDuration);
            return (int)duration.TotalSeconds;
        }
    }
}