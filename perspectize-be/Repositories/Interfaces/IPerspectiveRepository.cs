using perspectize_be.DTOs;

namespace perspectize_be.Repositories
{
    public interface IPerspectiveRepository
    {
        Task<IEnumerable<PerspectiveResponse>> GetPerspectivesByUsernameAsync(string username);
        Task<PerspectiveResponse?> GetPerspectiveByIdAsync(int id);
        Task<int> CreatePerspectivesAsync(IEnumerable<CreatePerspectiveRequest> requests);
        Task<int> UpdatePerspectiveAsync(int id, UpdatePerspectiveRequest request);
        Task<int> DeletePerspectivesAsync(int[] perspectiveIds);
        Task<bool> PerspectiveExistsAsync(int id);
    }
}