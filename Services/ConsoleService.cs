namespace _.Services;

public class ConsoleService
{
    private readonly ILogger<ConsoleService> _logger;

    public ConsoleService(ILogger<ConsoleService> logger)
    {
        _logger = logger;
    }

    public void PrintHello()
    {
        Console.WriteLine("hello");
        _logger.LogInformation("Hello job executed at: {time}", DateTimeOffset.Now);
    }
    public void UpdateAccountRegistry()
    {   
        _logger.LogInformation("Update account registry started at: {time}", DateTimeOffset.Now);
    }
} 