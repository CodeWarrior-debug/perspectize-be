// Controllers/ContentController.cs
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using perspectize_be.Data;
using perspectize_be.Models;

namespace perspectize_be.Controllers
{
    [Route("content")]
    [ApiController]
    public class ContentController : ControllerBase
    {
        private readonly ApplicationDbContext _context;

        public ContentController(ApplicationDbContext context)
        {
            _context = context;
        }

        // Endpoint 1: GET /content - gathers all
        [HttpGet]
        public async Task<IActionResult> GetAllContent(CancellationToken cancellationToken)
        {
            List<Content> contents = await _context.Contents
                .ToListAsync(cancellationToken);

            return Ok(contents.Select(c => new
                {
                    c.Id,
                    c.Name,
                    c.Url,
                    c.ContentType,
                    c.Length,
                    c.LengthUnits,
                    c.Response,
                    c.CreatedAt,
                    c.UpdatedAt
                }));
        }

        // Endpoint 2: GET /content/{name} - returns by unique name of content
        [HttpGet("{name}")]
        public async Task<IActionResult> GetContentByName(string name, CancellationToken cancellationToken)
        {
            if (string.IsNullOrEmpty(name))
            {
                return BadRequest(new { message = "Name parameter is required" });
            }

            Content? content = await _context.Contents
                .FirstOrDefaultAsync(c => c.Name == name, cancellationToken);

            if (content == null)
            {
                return NotFound(new { message = $"Content with name '{name}' not found" });
            }

            return Ok(new
            {
                content.Id,
                content.Name,
                content.Url,
                content.ContentType,
                content.Length,
                content.LengthUnits,
                content.Response,
                content.CreatedAt,
                content.UpdatedAt
            });
        }
    }
}