using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace perspectize_be.Migrations
{
    /// <inheritdoc />
    public partial class InitialCreate : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.Sql(@"
                CREATE TABLE public.content (
                    url varchar NULL,
                    length varchar NULL,
                    length_units varchar NULL,
                    response jsonb NULL,
                    content_type varchar NOT NULL,
                    name varchar NOT NULL,
                    id serial NOT NULL,
                    created_at timestamptz DEFAULT NOW() NOT NULL,
                    updated_at timestamptz DEFAULT NOW() NOT NULL,
                    CONSTRAINT content_unique_url UNIQUE(url),
                    CONSTRAINT content_pk PRIMARY KEY(id),
                    CONSTRAINT content_unique_name UNIQUE(name)
                );
            ");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.Sql(@"DROP TABLE public.content;");
        }
    }
}
