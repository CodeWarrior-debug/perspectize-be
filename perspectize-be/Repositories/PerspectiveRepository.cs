using Dapper;
using System.Data;
using System.Text.Json;
using perspectize_be.DTOs;

namespace perspectize_be.Repositories
{
    public class PerspectiveRepository : IPerspectiveRepository
    {
        private readonly IDbConnection _dbConnection;

        public PerspectiveRepository(IDbConnection dbConnection)
        {
            _dbConnection = dbConnection;
        }

        public async Task<IEnumerable<PerspectiveResponse>> GetPerspectivesByUsernameAsync(string username)
        {
            const string sql = @"
                SELECT 
                    c.id AS ContentId,
                    c.name AS Name,
                    c.content_type AS Type,
                    c.url AS Url,
                    p.id AS PerspectiveId,
                    p.like AS Like,
                    p.quality AS Quality,
                    p.agreement AS Agreement,
                    p.importance AS Importance,
                    p.privacy AS Privacy,
                    p.parts AS Parts,
                    p.created_at AS CreatedAt,
                    p.updated_at AS UpdatedAt,
                    p.categorized_ratings AS CategorizedRatings,
                    p.confidence AS Confidence,
                    p.user_id AS UserId,
                    p.claim AS Claim,
                    p.review_status AS ReviewStatus,
                    p.description AS Description,
                    p.labels AS Labels,
                    p.category AS Category
                FROM content c
                INNER JOIN perspectives p ON c.id = p.content_id
                INNER JOIN users u ON p.user_id = u.id
                WHERE u.username = @Username
                ORDER BY p.updated_at DESC";

            return await _dbConnection.QueryAsync<PerspectiveResponse>(sql, new { Username = username });
        }

        public async Task<PerspectiveResponse?> GetPerspectiveByIdAsync(int id)
        {
            const string sql = @"
                SELECT 
                    c.id AS ContentId,
                    c.name AS Name,
                    c.content_type AS Type,
                    c.url AS Url,
                    p.id AS PerspectiveId,
                    p.like AS Like,
                    p.quality AS Quality,
                    p.agreement AS Agreement,
                    p.importance AS Importance,
                    p.privacy AS Privacy,
                    p.parts AS Parts,
                    p.created_at AS CreatedAt,
                    p.updated_at AS UpdatedAt,
                    p.categorized_ratings AS CategorizedRatings,
                    p.confidence AS Confidence,
                    p.user_id AS UserId,
                    p.claim AS Claim,
                    p.review_status AS ReviewStatus,
                    p.description AS Description,
                    p.labels AS Labels,
                    p.category AS Category
                FROM content c
                INNER JOIN perspectives p ON c.id = p.content_id
                WHERE p.id = @Id";

            return await _dbConnection.QueryFirstOrDefaultAsync<PerspectiveResponse>(sql, new { Id = id });
        }

        public async Task<int> CreatePerspectivesAsync(IEnumerable<CreatePerspectiveRequest> requests)
        {
            const string sql = @"
                INSERT INTO perspectives (
                    like, quality, agreement, importance, privacy, parts,
                    categorized_ratings, confidence, user_id, claim, review_status,
                    description, labels, category, content_id
                ) VALUES (
                    @Like, @Quality, @Agreement, @Importance, @Privacy, @Parts,
                    @CategorizedRatings::jsonb[], @Confidence, @UserId, @Claim, @ReviewStatus,
                    @Description, @Labels, @Category, @ContentId
                )";

            IEnumerable<object> parameters = requests.Select(r => new
            {
                r.Like,
                r.Quality,
                r.Agreement,
                r.Importance,
                r.Privacy,
                r.Parts,
                CategorizedRatings = r.CategorizedRatings != null ? JsonSerializer.Serialize(r.CategorizedRatings) : null,
                r.Confidence,
                r.UserId,
                r.Claim,
                r.ReviewStatus,
                r.Description,
                r.Labels,
                r.Category,
                r.ContentId
            });

            return await _dbConnection.ExecuteAsync(sql, parameters);
        }

        public async Task<int> UpdatePerspectiveAsync(int id, UpdatePerspectiveRequest request)
        {
            List<string> setParts = new List<string>();
            DynamicParameters parameters = new DynamicParameters();
            parameters.Add("Id", id);

            if (request.Like != null)
            {
                setParts.Add("like = @Like");
                parameters.Add("Like", request.Like);
            }
            if (request.Quality.HasValue)
            {
                setParts.Add("quality = @Quality");
                parameters.Add("Quality", request.Quality);
            }
            if (request.Agreement.HasValue)
            {
                setParts.Add("agreement = @Agreement");
                parameters.Add("Agreement", request.Agreement);
            }
            if (request.Importance.HasValue)
            {
                setParts.Add("importance = @Importance");
                parameters.Add("Importance", request.Importance);
            }
            if (request.Privacy != null)
            {
                setParts.Add("privacy = @Privacy");
                parameters.Add("Privacy", request.Privacy);
            }
            if (request.Parts != null)
            {
                setParts.Add("parts = @Parts");
                parameters.Add("Parts", request.Parts);
            }
            if (request.CategorizedRatings != null)
            {
                setParts.Add("categorized_ratings = @CategorizedRatings::jsonb[]");
                parameters.Add("CategorizedRatings", JsonSerializer.Serialize(request.CategorizedRatings));
            }
            if (request.Confidence.HasValue)
            {
                setParts.Add("confidence = @Confidence");
                parameters.Add("Confidence", request.Confidence);
            }
            if (request.Claim != null)
            {
                setParts.Add("claim = @Claim");
                parameters.Add("Claim", request.Claim);
            }
            if (request.ReviewStatus != null)
            {
                setParts.Add("review_status = @ReviewStatus");
                parameters.Add("ReviewStatus", request.ReviewStatus);
            }
            if (request.Description != null)
            {
                setParts.Add("description = @Description");
                parameters.Add("Description", request.Description);
            }
            if (request.Labels != null)
            {
                setParts.Add("labels = @Labels");
                parameters.Add("Labels", request.Labels);
            }
            if (request.Category != null)
            {
                setParts.Add("category = @Category");
                parameters.Add("Category", request.Category);
            }

            if (!setParts.Any())
            {
                return 0; // No fields to update
            }

            string updateSql = $@"
                UPDATE perspectives 
                SET {string.Join(", ", setParts)}
                WHERE id = @Id";

            return await _dbConnection.ExecuteAsync(updateSql, parameters);
        }

        public async Task<int> DeletePerspectivesAsync(int[] perspectiveIds)
        {
            const string sql = "DELETE FROM perspectives WHERE id = ANY(@Ids)";
            return await _dbConnection.ExecuteAsync(sql, new { Ids = perspectiveIds });
        }

        public async Task<bool> PerspectiveExistsAsync(int id)
        {
            const string sql = "SELECT 1 FROM perspectives WHERE id = @Id LIMIT 1";
            int? result = await _dbConnection.QueryFirstOrDefaultAsync<int?>(sql, new { Id = id });
            return result.HasValue;
        }
    }
}