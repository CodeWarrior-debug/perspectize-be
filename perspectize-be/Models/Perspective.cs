using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace perspectize_be.Models
{
    [Table("perspectives")]
    public class Perspective
    {
        [Key]
        [Column("id")]
        public int Id { get; set; }
        
        [Column("like")]
        public string? Like { get; set; }
        
        [Column("quality")]
        public int? Quality { get; set; }
        
        [Column("agreement")]
        public int? Agreement { get; set; }
        
        [Column("importance")]
        public int? Importance { get; set; }
        
        [Column("privacy")]
        public string? Privacy { get; set; }
        
        [Column("parts")]
        public int[]? Parts { get; set; }
        
        [Column("created_at")]
        public DateTime CreatedAt { get; set; }
        
        [Column("updated_at")]
        public DateTime UpdatedAt { get; set; }
        
        [Column("categorized_ratings", TypeName = "jsonb")]
        public object[]? CategorizedRatings { get; set; }
        
        [Column("confidence")]
        public int? Confidence { get; set; }
        
        [Column("user_id")]
        public int UserId { get; set; }
        
        [Required]
        [MaxLength(255)]
        [Column("claim")]
        public string Claim { get; set; } = string.Empty;
        
        [Column("review_status")]
        public string? ReviewStatus { get; set; }
        
        [Column("description")]
        public string? Description { get; set; }
        
        [Column("labels")]
        public string[]? Labels { get; set; }
        
        [Column("category")]
        public string? Category { get; set; }
        
        [Column("content_id")]
        public int? ContentId { get; set; }
    }
}