using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace perspectize_be.Migrations
{
    /// <inheritdoc />
    public partial class UpdateLengthToNumericType : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            // Use raw SQL to convert the column type with USING clause
            migrationBuilder.Sql("ALTER TABLE content ALTER COLUMN length TYPE integer USING length::integer;");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            // Convert back to varchar in down migration
            migrationBuilder.Sql("ALTER TABLE content ALTER COLUMN length TYPE varchar USING length::varchar;");
        }
    }
}
