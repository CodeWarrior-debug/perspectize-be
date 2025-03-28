using Hangfire;
using Microsoft.Extensions.Configuration;

namespace _.Services;

public class CronJobsService
{
    private readonly IConfiguration _configuration;

    public CronJobsService(IConfiguration configuration)
    {
        _configuration = configuration;
    }

    public void RegisterJobs()
    {
        string registryUpdateSchedule = _configuration.GetValue<string>("Logging:Configs:REGISTRY_UPDATE_SCHEDULE");
        
        RecurringJob.AddOrUpdate<ConsoleService>(
            "print-hello-job",
            service => service.PrintHello(),
            Cron.Minutely());

        RecurringJob.AddOrUpdate<ConsoleService>(
            "update-account-registry",
            service => service.UpdateAccountRegistry(),
            registryUpdateSchedule);
    }
} 