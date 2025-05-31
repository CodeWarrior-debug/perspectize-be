using Microsoft.EntityFrameworkCore;
using perspectize_be.Data;
using Npgsql.EntityFrameworkCore.PostgreSQL;
using perspectize_be.Services;
using System.Data;
using Npgsql;
using perspectize_be.Repositories;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers()
    .AddJsonOptions(options =>
    {
        options.JsonSerializerOptions.PropertyNamingPolicy = System.Text.Json.JsonNamingPolicy.CamelCase;
    });

// Configure DbContext
builder.Services.AddDbContext<ApplicationDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));

// Register IDbConnection for Dapper
builder.Services.AddScoped<IDbConnection>(sp => 
    new NpgsqlConnection(builder.Configuration.GetConnectionString("DefaultConnection")));

// Register HttpClient
builder.Services.AddHttpClient();

// Register repositories
builder.Services.AddScoped<IPerspectiveRepository, PerspectiveRepository>();

// Register YouTube service
builder.Services.AddScoped<YouTubeService>();
builder.Services.AddScoped<IPerspectiveService, PerspectiveService>();


var app = builder.Build();

// Configure the HTTP request pipeline
if (app.Environment.IsDevelopment())
{
    try
    {
        // Apply migrations and seed data in development
        using (IServiceScope scope = app.Services.CreateScope())
        {
            ApplicationDbContext db = scope.ServiceProvider.GetRequiredService<ApplicationDbContext>();
            db.Database.Migrate();

            // Seed data
            await SeedData.InitializeAsync(app.Services);
        }
    }
    catch (Exception ex)
    {
        ILogger<Program> logger = app.Services.GetRequiredService<ILogger<Program>>();
        logger.LogError(ex, "An error occurred while migrating or seeding the database.");
    }
}

// TODO: stop localhost 7253 opening browswer window every time

app.UseRouting();

app.MapControllers();

app.Run();