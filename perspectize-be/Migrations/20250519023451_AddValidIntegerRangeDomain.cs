using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;
using System;

#nullable disable

namespace perspectize_be.Migrations
{
    /// <inheritdoc />
    public partial class AddValidIntegerRangeDomain : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.Sql(@"
                DO $$ 
                BEGIN
                    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'valid_integer_range') THEN
                        CREATE DOMAIN valid_integer_range AS integer
                        CHECK (VALUE BETWEEN 0 AND 10000);
                    END IF;
                END $$;
            ");

            migrationBuilder.CreateTable(
                name: "users",
                schema: "public",
                columns: table => new
                {
                    id = table.Column<int>(type: "integer", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    username = table.Column<string>(type: "character varying(24)", maxLength: 24, nullable: false),
                    email = table.Column<string>(type: "text", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("users_pk", x => x.id);
                    table.UniqueConstraint("users_unique_email", x => x.email);
                    table.UniqueConstraint("users_unique_username", x => x.username);
                });

            migrationBuilder.CreateTable(
                name: "perspectives",
                schema: "public",
                columns: table => new
                {
                    id = table.Column<int>(type: "integer", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    like = table.Column<string>(type: "text", nullable: true),
                    quality = table.Column<int>(type: "valid_integer_range", nullable: true),
                    agreement = table.Column<int>(type: "valid_integer_range", nullable: true),
                    importance = table.Column<int>(type: "valid_integer_range", nullable: true),
                    privacy = table.Column<string>(type: "text", nullable: true, defaultValue: "public"),
                    parts = table.Column<int[]>(type: "integer[]", nullable: true),
                    created_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false, defaultValueSql: "now()"),
                    updated_at = table.Column<DateTime>(type: "timestamp with time zone", nullable: false, defaultValueSql: "now()"),
                    categorized_ratings = table.Column<object[]>(type: "jsonb[]", nullable: true),
                    confidence = table.Column<int>(type: "valid_integer_range", nullable: true),
                    user_id = table.Column<int>(type: "integer", nullable: false),
                    claim = table.Column<string>(type: "character varying(255)", maxLength: 255, nullable: false),
                    review_status = table.Column<string>(type: "text", nullable: true),
                    description = table.Column<string>(type: "text", nullable: true),
                    labels = table.Column<string[]>(type: "text[]", nullable: true),
                    category = table.Column<string>(type: "text", nullable: true),
                    content_id = table.Column<int>(type: "integer", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("perspectives_pk", x => x.id);
                    table.UniqueConstraint("perspectives_unique_user_claims", x => new { x.claim, x.user_id });
                    table.ForeignKey(
                        name: "perspectives_content_fk",
                        column: x => x.content_id,
                        principalSchema: "public",
                        principalTable: "content",
                        principalColumn: "id");
                    table.ForeignKey(
                        name: "perspectives_users_fk",
                        column: x => x.user_id,
                        principalSchema: "public",
                        principalTable: "users",
                        principalColumn: "id");
                });

            migrationBuilder.Sql(@"
                CREATE OR REPLACE FUNCTION public.update_updated_at()
                RETURNS trigger
                LANGUAGE plpgsql
                AS $function$
                    BEGIN
                        NEW.updated_at = now();
                        RETURN NEW;
                    END;
                $function$;
            ");

            migrationBuilder.Sql(@"
                CREATE TRIGGER set_updated_at
                    BEFORE UPDATE ON public.perspectives
                    FOR EACH ROW
                    EXECUTE FUNCTION update_updated_at();
            ");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.Sql(@"
                DROP TRIGGER IF EXISTS set_updated_at ON public.perspectives;
            ");

            migrationBuilder.DropTable(
                name: "perspectives",
                schema: "public");

            migrationBuilder.DropTable(
                name: "users",
                schema: "public");

            migrationBuilder.Sql(@"
                DROP FUNCTION IF EXISTS public.update_updated_at();
            ");

            migrationBuilder.Sql(@"
                DO $$ 
                BEGIN
                    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'valid_integer_range') THEN
                        DROP DOMAIN IF EXISTS valid_integer_range;
                    END IF;
                END $$;
            ");
        }
    }
} 