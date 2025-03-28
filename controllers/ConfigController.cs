using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;

namespace _.Controllers;

[ApiController]
[Route("[controller]")]
public class ConfigController : ControllerBase
{
    private readonly IConfiguration _configuration;

    public ConfigController(IConfiguration configuration)
    {
        _configuration = configuration;
    }

    [HttpGet("refill-remain")]
    public IActionResult GetRefillRemain()
    {
        // Direct access to the configuration value
        int numberRemain = _configuration.GetValue<int>("Logging:Configs:REFILL_ACCTIDS_NUMBER_REMAIN");
        
        return Ok(new { RefillAccountIdsNumberRemain = numberRemain });
    }
} 