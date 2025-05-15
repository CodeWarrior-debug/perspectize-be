using System.Text.Json.Serialization;

namespace perspectize_be.DTOs
{
    public class VideoRequest
    {
        public string VideoId { get; set; } = string.Empty;
    }
    
    public class VideosRequest
    {
        public List<string> VideoUrls { get; set; } = new List<string>();
    }
}