using Microsoft.EntityFrameworkCore;
using perspectize_be.Data;
using Npgsql.EntityFrameworkCore.PostgreSQL;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();

// Configure DbContext
builder.Services.AddDbContext<ApplicationDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));

var app = builder.Build();

// TODO: stop localhost 7253 opening browswer window every time

app.UseRouting();

app.MapControllers();

app.Run();