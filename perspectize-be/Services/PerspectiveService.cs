using perspectize_be.DTOs;
using perspectize_be.Repositories;

namespace perspectize_be.Services
{
    public class PerspectiveService : IPerspectiveService
    {
        private readonly IPerspectiveRepository _perspectiveRepository;

        public PerspectiveService(IPerspectiveRepository perspectiveRepository)
        {
            _perspectiveRepository = perspectiveRepository;
        }

        public async Task<IEnumerable<PerspectiveResponse>> GetPerspectivesByUsernameAsync(string username)
        {
            return await _perspectiveRepository.GetPerspectivesByUsernameAsync(username);
        }

        public async Task<PerspectiveResponse?> GetPerspectiveByIdAsync(int id)
        {
            return await _perspectiveRepository.GetPerspectiveByIdAsync(id);
        }

        public async Task<int> CreatePerspectivesAsync(IEnumerable<CreatePerspectiveRequest> requests)
        {
            try
            {
                return await _perspectiveRepository.CreatePerspectivesAsync(requests);
            }
            catch (Exception ex) when (ex.Message.Contains("perspectives_unique_user_claims"))
            {
                throw new InvalidOperationException("A perspective with this claim already exists for the user.");
            }
            catch (Exception ex) when (ex.Message.Contains("valid_integer_range"))
            {
                throw new ArgumentOutOfRangeException("Quality, agreement, importance, and confidence values must be between 0 and 10000.");
            }
        }

        public async Task<PerspectiveResponse?> UpdatePerspectiveAsync(int id, UpdatePerspectiveRequest request)
        {
            try
            {
                // Check if perspective exists first
                bool exists = await _perspectiveRepository.PerspectiveExistsAsync(id);
                if (!exists)
                    return null;

                int affectedRows = await _perspectiveRepository.UpdatePerspectiveAsync(id, request);
                if (affectedRows == 0)
                    return null;

                // Return updated perspective
                return await _perspectiveRepository.GetPerspectiveByIdAsync(id);
            }
            catch (Exception ex) when (ex.Message.Contains("perspectives_unique_user_claims"))
            {
                throw new InvalidOperationException("A perspective with this claim already exists for the user.");
            }
            catch (Exception ex) when (ex.Message.Contains("valid_integer_range"))
            {
                throw new ArgumentOutOfRangeException("Quality, agreement, importance, and confidence values must be between 0 and 10000.");
            }
        }

        public async Task<int> DeletePerspectivesAsync(int[] perspectiveIds)
        {
            return await _perspectiveRepository.DeletePerspectivesAsync(perspectiveIds);
        }
    }
}