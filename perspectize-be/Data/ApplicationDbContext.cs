using Microsoft.EntityFrameworkCore;
using perspectize_be.Models;

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

         public DbSet<Content> Contents { get; set; }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            base.OnModelCreating(modelBuilder);
            
            // Configure Content entity
            modelBuilder.Entity<Content>()
                .HasIndex(c => c.Url)
                .IsUnique();
                
            modelBuilder.Entity<Content>()
                .HasIndex(c => c.Name)
                .IsUnique();
        }
    }
} 