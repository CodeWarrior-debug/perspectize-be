using Microsoft.EntityFrameworkCore;

namespace perspectize_be.Data
{
    public class ApplicationDbContext : DbContext
    {
        public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
            : base(options)
        {
        }

        // Add DbSet properties for your entities here
        // Example:
        // public DbSet<Video> Videos { get; set; }
    }
} 