namespace perspectize_be.DTOs
{
    public class PerspectiveResponse
    {
        public int ContentId { get; set; }
        public string Name { get; set; } = string.Empty;
        public string Type { get; set; } = string.Empty; 
        public string? Url { get; set; }
        
        // Perspective fields
        public int PerspectiveId { get; set; }
        public string? Like { get; set; }
        public int? Quality { get; set; }
        public int? Agreement { get; set; }
        public int? Importance { get; set; }
        public string? Privacy { get; set; }
        public int[]? Parts { get; set; }
        public DateTime CreatedAt { get; set; }
        public DateTime UpdatedAt { get; set; }
        public object[]? CategorizedRatings { get; set; }
        public int? Confidence { get; set; }
        public int UserId { get; set; }
        public string Claim { get; set; } = string.Empty;
        public string? ReviewStatus { get; set; }
        public string? Description { get; set; }
        public string[]? Labels { get; set; }
        public string? Category { get; set; }
    }
}