using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using System.Threading;
using Microsoft.AspNetCore.Authorization;
using App.DTOs;

namespace App.Controllers
{
    [Route("youtube")]
    [ApiController]
    [AllowAnonymous]
    public class YTController : ControllerBase
    {
        [HttpPost("video")]
        public IActionResult PostVideo([FromBody] VideoRequest req, CancellationToken cancellationToken)
        {
            if (string.IsNullOrEmpty(req.VideoId)) return BadRequest(new { message = "videoId is required" });
            return Ok(
                new{message = $"Video '{req.VideoId}' posted successfully"}
            );
        }
    }
}
