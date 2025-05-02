using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using System.Text.Json.Serialization;

namespace perspectize_be.Models
{
    [Table("content")]
    public class Content
    {
        [Key]
        [Column("id")]
        public int Id { get; set; }

        [Column("url")]
        public string? Url { get; set; }

        [Column("length")]
        public string? Length { get; set; }

        [Column("length_units")]
        public string? LengthUnits { get; set; }

        [Column("response")]
        [JsonIgnore]
        public string? ResponseJson { get; set; }

        [NotMapped]
        public object? Response 
        {
            get => !string.IsNullOrEmpty(ResponseJson) 
                ? System.Text.Json.JsonSerializer.Deserialize<object>(ResponseJson) 
                : null;
            set => ResponseJson = value != null 
                ? System.Text.Json.JsonSerializer.Serialize(value) 
                : null;
        }

        [Column("content_type")]
        [Required]
        public string ContentType { get; set; } = "youtube";

        [Column("name")]
        [Required]
        public string Name { get; set; } = string.Empty;

        [Column("created_at")]
        public DateTime CreatedAt { get; set; }

        [Column("updated_at")]
        public DateTime UpdatedAt { get; set; }
    }
}