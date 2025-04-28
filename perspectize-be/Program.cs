var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();

var app = builder.Build();

// TODO: stop localhost 7253 opening browswer window every time

app.UseRouting();

app.MapControllers();

app.Run();