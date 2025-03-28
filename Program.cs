using Hangfire;
using Hangfire.MemoryStorage;
using _.Services;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

builder.Services.AddScoped<ConsoleService>();
builder.Services.AddSingleton<CronJobsService>();

builder.Services.AddHangfire(config =>
{
    config.UseMemoryStorage();
});
builder.Services.AddHangfireServer();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.UseAuthorization();
app.UseHangfireDashboard();

app.MapControllers();

app.Services.GetRequiredService<CronJobsService>().RegisterJobs();

app.Run();
