using Microsoft.AspNetCore.Mvc;
using Dapper;
using System.Data;
using perspectize_be.Models;

namespace perspectize_be.Controllers
{
    [Route("content")]
    [ApiController]
    public class ContentController : ControllerBase
    {
        private readonly IDbConnection _dbConnection;

        public ContentController(IDbConnection dbConnection)
        {
            _dbConnection = dbConnection;
        }

        [HttpGet]
        public async Task<IActionResult> GetAllContent(CancellationToken cancellationToken)
        {
            string query = @"
                SELECT id AS Id, name AS Name, url AS Url, content_type AS ContentType, 
                       length AS Length, length_units AS LengthUnits, response AS Response,
                       created_at AS CreatedAt, updated_at AS UpdatedAt
                FROM content";
            
            IEnumerable<Content> contents = await _dbConnection.QueryAsync<Content>(query);
            return Ok(contents);
        }

        [HttpGet("{name}")] //TODO: later, change to id, names include spaces and can get long
        public async Task<IActionResult> GetContentByName(string name, CancellationToken cancellationToken)
        {
            if (string.IsNullOrEmpty(name))
            {
                return BadRequest(new { message = "Name parameter is required" });
            }

            string query = @"
                SELECT id AS Id, name AS Name, url AS Url, content_type AS ContentType, 
                       length AS Length, length_units AS LengthUnits, response AS Response,
                       created_at AS CreatedAt, updated_at AS UpdatedAt
                FROM content
                WHERE name = @Name
                LIMIT 1";
            
            Content? content = await _dbConnection.QueryFirstOrDefaultAsync<Content>(query, new { Name = name });

            if (content == null)
            {
                return NotFound(new { message = $"Content with name '{name}' not found" });
            }

            return Ok(content);
        }
    }
}