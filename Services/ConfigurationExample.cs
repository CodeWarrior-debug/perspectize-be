using Microsoft.Extensions.Configuration;
using System;

namespace _.Services;

public class ConfigurationExample
{
    private readonly IConfiguration _configuration;

    public ConfigurationExample(IConfiguration configuration)
    {
        _configuration = configuration;
    }

    public int GetRefillAccountIdsNumberRemain()
    {
        int numberRemain = _configuration.GetValue<int>("Logging:Configs:REFILL_ACCTIDS_NUMBER_REMAIN");
        int percentUsed = _configuration.GetValue<int>("Logging:Configs:REFILL_ACCTIDS_PERCENT_USED");
        int numberToAdd = _configuration.GetValue<int>("Logging:Configs:REFILL_ACCTIDS_NUMBER_TO_ADD");
        string refillCheckSchedule = _configuration.GetValue<string>("Logging:Configs:REFILL_CHECK_SCHEDULE");
        
        return numberRemain;
    }

    // Alternative approach using configuration binding
    public void DemonstrateConfigurationBinding()
    {
        var configsSection = _configuration.GetSection("Logging:Configs");
        string refillCheckSchedule = configsSection["REFILL_CHECK_SCHEDULE"];
        int numberRemain = int.Parse(configsSection["REFILL_ACCTIDS_NUMBER_REMAIN"]);
        int percentUsed = int.Parse(configsSection["REFILL_ACCTIDS_PERCENT_USED"]);
        int numberToAdd = int.Parse(configsSection["REFILL_ACCTIDS_NUMBER_TO_ADD"]);

        
        Console.WriteLine($"REFILL_ACCTIDS_NUMBER_REMAIN value: {numberRemain}");
    }
} 