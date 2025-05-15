using System;
using System.Text.Json;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

#nullable disable

namespace perspectize_be.Migrations
{
    /// <inheritdoc />
    public partial class UpdateResponseColumnToJsonbType : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            // Use raw SQL to modify the column type to jsonb
            migrationBuilder.Sql(@"
                DO $$
                BEGIN
                    -- First check if the column exists
                    IF EXISTS (
                        SELECT FROM information_schema.columns 
                        WHERE table_schema = 'public' 
                        AND table_name = 'content' 
                        AND column_name = 'response'
                    ) THEN
                        -- If response column exists and isn't jsonb, convert it
                        IF (SELECT data_type FROM information_schema.columns 
                            WHERE table_schema = 'public' 
                            AND table_name = 'content' 
                            AND column_name = 'response') != 'jsonb' THEN
                            -- Try to convert existing data to jsonb
                            BEGIN
                                ALTER TABLE content 
                                ALTER COLUMN response TYPE jsonb USING response::jsonb;
                            EXCEPTION WHEN OTHERS THEN
                                -- If conversion fails, drop and recreate the column
                                ALTER TABLE content DROP COLUMN response;
                                ALTER TABLE content ADD COLUMN response jsonb;
                            END;
                        END IF;
                    ELSE
                        -- If response column doesn't exist, add it
                        ALTER TABLE content ADD COLUMN response jsonb;
                    END IF;
                END $$;
            ");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            // Revert to text type if needed
            migrationBuilder.Sql(@"
                ALTER TABLE content 
                ALTER COLUMN response TYPE text USING response::text;
            ");
        }
    }
}
