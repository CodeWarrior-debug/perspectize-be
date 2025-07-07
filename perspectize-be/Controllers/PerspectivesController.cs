using Microsoft.AspNetCore.Mvc;
using perspectize_be.DTOs;
using perspectize_be.Services;

namespace perspectize_be.Controllers
{
    [ApiController]
    [Route("perspectives")]
    public class PerspectivesController : ControllerBase
    {
        private readonly IPerspectiveService _perspectiveService;

        public PerspectivesController(IPerspectiveService perspectiveService)
        {
            _perspectiveService = perspectiveService;
        }

        [HttpGet("{username}")]
        public async Task<ActionResult<IEnumerable<PerspectiveResponse>>> GetPerspectivesByUsername(string username)
        {
            IEnumerable<PerspectiveResponse> perspectives = await _perspectiveService.GetPerspectivesByUsernameAsync(username);
            return Ok(perspectives);
        }

        [HttpGet("{id:int}")]
        public async Task<ActionResult<PerspectiveResponse>> GetPerspectiveById(int id)
        {
            PerspectiveResponse? perspective = await _perspectiveService.GetPerspectiveByIdAsync(id);
            if (perspective == null)
                return NotFound();
            
            return Ok(perspective);
        }

        [HttpPost]
        public async Task<ActionResult> CreatePerspectives([FromBody] List<CreatePerspectiveRequest> requests)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);
                
            if (!requests.Any())
                return BadRequest("At least one perspective is required");

            try
            {
                int createdCount = await _perspectiveService.CreatePerspectivesAsync(requests);
                return StatusCode(201, new { created_count = createdCount });
            }
            catch (InvalidOperationException ex)
            {
                return Conflict(new { message = ex.Message });
            }
            catch (ArgumentOutOfRangeException ex)
            {
                return BadRequest(new { message = ex.Message });
            }
        }

        [HttpPut("{id}")]
        public async Task<ActionResult<PerspectiveResponse>> UpdatePerspective(int id, [FromBody] UpdatePerspectiveRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            try
            {
                PerspectiveResponse? updated = await _perspectiveService.UpdatePerspectiveAsync(id, request);
                if (updated == null)
                    return NotFound();
                
                return Ok(updated);
            }
            catch (InvalidOperationException ex)
            {
                return Conflict(new { message = ex.Message });
            }
            catch (ArgumentOutOfRangeException ex)
            {
                return BadRequest(new { message = ex.Message });
            }
        }

        [HttpDelete]
        public async Task<ActionResult> DeletePerspectives([FromBody] int[] perspectiveIds)
        {
            if (!perspectiveIds.Any())
                return BadRequest("At least one perspective ID is required");

            int deletedCount = await _perspectiveService.DeletePerspectivesAsync(perspectiveIds);
            return StatusCode(204, new { deleted_count = deletedCount });
        }
    }
}