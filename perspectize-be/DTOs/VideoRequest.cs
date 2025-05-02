using System.Text.Json.Serialization;

namespace perspectize_be.DTOs
{
    // Request DTOs
    public class VideoRequest
    {
        public string VideoId { get; set; } = string.Empty;
    }
    namespace perspectize_be.DTOs
{
    public class VideosRequest
    {
        public List<string> VideoUrls { get; set; } = new List<string>();
    }
}

}