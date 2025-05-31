using perspectize_be.DTOs;

namespace perspectize_be.Services
{
    public interface IPerspectiveService
    {
        Task<IEnumerable<PerspectiveResponse>> GetPerspectivesByUsernameAsync(string username);
        Task<PerspectiveResponse?> GetPerspectiveByIdAsync(int id);
        Task<int> CreatePerspectivesAsync(IEnumerable<CreatePerspectiveRequest> requests);
        Task<PerspectiveResponse?> UpdatePerspectiveAsync(int id, UpdatePerspectiveRequest request);
        Task<int> DeletePerspectivesAsync(int[] perspectiveIds);
    }
}