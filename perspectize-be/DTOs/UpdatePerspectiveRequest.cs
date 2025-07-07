using System.ComponentModel.DataAnnotations;

namespace perspectize_be.DTOs
{
    public class UpdatePerspectiveRequest
    {
        public string? Like { get; set; }
        
        [Range(0, 10000, ErrorMessage = "Quality must be between 0 and 10000")]
        public int? Quality { get; set; }
        
        [Range(0, 10000, ErrorMessage = "Agreement must be between 0 and 10000")]
        public int? Agreement { get; set; }
        
        [Range(0, 10000, ErrorMessage = "Importance must be between 0 and 10000")]
        public int? Importance { get; set; }
        
        public string? Privacy { get; set; }
        public int[]? Parts { get; set; }
        public object[]? CategorizedRatings { get; set; }
        
        [Range(0, 10000, ErrorMessage = "Confidence must be between 0 and 10000")]
        public int? Confidence { get; set; }
        
        [MaxLength(255)]
        public string? Claim { get; set; }
        
        public string? ReviewStatus { get; set; }
        public string? Description { get; set; }
        public string[]? Labels { get; set; }
        public string? Category { get; set; }
    }
}