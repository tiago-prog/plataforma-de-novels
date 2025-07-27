-- BEGIN;

-- -- 0. Enable required PostgreSQL extensions
-- CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- -- 1. Criação de ENUMs
-- CREATE TYPE user_role AS ENUM ('reader', 'author', 'moderator', 'admin');
-- CREATE TYPE user_status AS ENUM ('active', 'suspended', 'banned', 'pending_verification');
-- CREATE TYPE work_type AS ENUM ('novel', 'manhwa', 'webtoon');
-- CREATE TYPE work_status AS ENUM ('draft', 'published', 'completed', 'hiatus', 'cancelled');
-- CREATE TYPE moderation_status AS ENUM ('pending', 'approved', 'rejected', 'flagged');
-- CREATE TYPE publication_status AS ENUM ('draft', 'scheduled', 'published', 'archived');
-- CREATE TYPE comment_status AS ENUM ('visible', 'hidden', 'deleted', 'flagged');
-- CREATE TYPE notification_type AS ENUM ('new_chapter', 'new_comment', 'new_follower', 'like_received', 'mention');      
-- CREATE TYPE report_status AS ENUM ('open', 'investigating', 'resolved', 'dismissed');
-- CREATE TYPE task_status AS ENUM ('pending', 'running', 'completed', 'failed');
-- CREATE TYPE task_type AS ENUM ('publish_content', 'send_notification', 'cleanup_data', 'backup_data');
-- CREATE TYPE target_type AS ENUM ('work', 'content', 'comment');

-- -- 2. Criação de Funções PRIMEIRO
-- CREATE OR REPLACE FUNCTION update_timestamp_column()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.updated_at = CURRENT_TIMESTAMP;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE OR REPLACE FUNCTION update_content_published_at()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     IF NEW.publication_status = 'published' 
--     AND OLD.publication_status != 'published' 
--     AND NEW.published_at IS NULL 
--     THEN
--         NEW.published_at = CURRENT_TIMESTAMP;
--     END IF;
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- -- 3. Criação de Tabelas
-- CREATE TABLE IF NOT EXISTS users (
--     id SERIAL PRIMARY KEY,
--     username VARCHAR(50) UNIQUE NOT NULL,
--     email VARCHAR(255) UNIQUE NOT NULL,
--     password_hash TEXT NOT NULL,
--     role user_role NOT NULL,
--     avatar TEXT,
--     verified BOOLEAN DEFAULT false,
--     status user_status DEFAULT 'pending_verification',
--     last_login_at TIMESTAMPTZ,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );

-- CREATE TABLE IF NOT EXISTS user_works (
--     id SERIAL PRIMARY KEY,
--     user_id INTEGER NOT NULL REFERENCES users(id),
--     work_id INTEGER NOT NULL REFERENCES works(id),
--     created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     UNIQUE(user_id, work_id)
-- );

-- CREATE TABLE IF NOT EXISTS categories (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(100) UNIQUE NOT NULL,
--     slug VARCHAR(120) UNIQUE NOT NULL,
--     description TEXT
-- );

-- CREATE TABLE IF NOT EXISTS genres (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(100) UNIQUE NOT NULL,
--     slug VARCHAR(120) UNIQUE NOT NULL,
--     description TEXT
-- );

-- CREATE TABLE IF NOT EXISTS works (
--     id SERIAL PRIMARY KEY,
--     title VARCHAR(255) NOT NULL,
--     slug VARCHAR(255) UNIQUE NOT NULL,
--     description TEXT,
--     author_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
--     type work_type NOT NULL,
--     work_status work_status DEFAULT 'draft',
--     moderation_status moderation_status DEFAULT 'pending',
--     cover_url TEXT,
--     published_at TIMESTAMPTZ,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );

-- -- Conteúdos (capítulos/episódios)
-- CREATE TABLE IF NOT EXISTS contents (
--     id SERIAL PRIMARY KEY,
--     work_id INT NOT NULL REFERENCES works(id) ON DELETE CASCADE,
--     number INT NOT NULL,
--     title VARCHAR(255) NOT NULL,
--     markdown_text TEXT,
--     display_order INT NOT NULL DEFAULT 0,
--     publication_status publication_status DEFAULT 'draft',
--     scheduled_for TIMESTAMPTZ,
--     published_at TIMestamptz,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     UNIQUE (work_id, number)
-- );

-- -- Imagens do conteúdo
-- CREATE TABLE IF NOT EXISTS content_images (
--     id SERIAL PRIMARY KEY,
--     content_id INT NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
--     image_url TEXT NOT NULL,
--     order_num INT NOT NULL DEFAULT 0,
--     alt_text VARCHAR(255)
-- );

-- -- Relação obras-categorias
-- CREATE TABLE IF NOT EXISTS work_categories (
--     work_id INT NOT NULL REFERENCES works(id) ON DELETE CASCADE,
--     category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
--     PRIMARY KEY (work_id, category_id)
-- );

-- -- Relação obras-gêneros
-- CREATE TABLE IF NOT EXISTS work_genres (
--     work_id INT NOT NULL REFERENCES works(id) ON DELETE CASCADE,
--     genre_id INT NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
--     PRIMARY KEY (work_id, genre_id)
-- );

-- -- Comentários
-- CREATE TABLE IF NOT EXISTS comments (
--     id SERIAL PRIMARY KEY,
--     user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     content_id INT NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
--     parent_id INT REFERENCES comments(id) ON DELETE CASCADE,
--     text_content TEXT NOT NULL,
--     comment_status comment_status DEFAULT 'visible',
--     positive_votes INT DEFAULT 0,
--     negative_votes INT DEFAULT 0,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );

-- -- Curtidas (sistema polimórfico)
-- CREATE TABLE IF NOT EXISTS likes (
--     id SERIAL PRIMARY KEY,
--     user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     target_type target_type NOT NULL,
--     target_id INT NOT NULL,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     UNIQUE (user_id, target_type, target_id)
-- );

-- -- Seguidores
-- CREATE TABLE IF NOT EXISTS followers (
--     follower_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     followed_user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     PRIMARY KEY (follower_user_id, followed_user_id),
--     CHECK (follower_user_id <> followed_user_id)
-- );

-- -- Notificações
-- CREATE TABLE IF NOT EXISTS notifications (
--     id SERIAL PRIMARY KEY,
--     user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--     notification_type notification_type NOT NULL,
--     payload JSONB,
--     read_at TIMESTAMPTZ,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );

-- -- Fila de moderação
-- CREATE TABLE IF NOT EXISTS moderation_queue (
--     id SERIAL PRIMARY KEY,
--     item_type target_type NOT NULL,
--     item_id INT NOT NULL,
--     submitted_by INT REFERENCES users(id) ON DELETE SET NULL,
--     moderation_status moderation_status DEFAULT 'pending',
--     justification TEXT,
--     reviewed_by INT REFERENCES users(id) ON DELETE SET NULL,
--     reviewed_at TIMESTAMPTZ,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );

-- -- Denúncias
-- CREATE TABLE IF NOT EXISTS reports (
--     id SERIAL PRIMARY KEY,
--     reporter_id INT REFERENCES users(id) ON DELETE SET NULL,
--     item_type VARCHAR(20) NOT NULL, -- work, content, comment, user
--     item_id INT NOT NULL,
--     reason TEXT,
--     report_status report_status DEFAULT 'open',
--     moderator_id INT REFERENCES users(id) ON DELETE SET NULL,
--     resolved_at TIMESTAMPTZ,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );

-- -- Tarefas agendadas
-- CREATE TABLE IF NOT EXISTS scheduled_tasks (
--     id SERIAL PRIMARY KEY,
--     task_type task_type NOT NULL,
--     reference_id INT, -- ID de referência opcional
--     execute_at TIMESTAMPTZ NOT NULL,
--     payload JSONB,
--     status task_status DEFAULT 'pending',
--     attempts INT DEFAULT 0,
--     last_error TEXT,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );

-- -- Log de auditoria
-- CREATE TABLE IF NOT EXISTS audit_log (
--     id SERIAL PRIMARY KEY,
--     entity VARCHAR(50) NOT NULL, -- Nome da entidade/tabela
--     entity_id INT NOT NULL,      -- ID do registro
--     action VARCHAR(10) NOT NULL, -- INSERT, UPDATE, DELETE
--     diff JSONB,                  -- Diferencial de alteração
--     performing_user_id INT REFERENCES users(id) ON DELETE SET NULL,
--     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );

-- -- 4. Criação de Triggers APÓS as funções e tabelas
-- CREATE TRIGGER users_update_timestamp
-- BEFORE UPDATE ON users
-- FOR EACH ROW EXECUTE FUNCTION update_timestamp_column();

-- CREATE TRIGGER works_update_timestamp
-- BEFORE UPDATE ON works
-- FOR EACH ROW EXECUTE FUNCTION update_timestamp_column();

-- CREATE TRIGGER contents_update_timestamp
-- BEFORE UPDATE ON contents
-- FOR EACH ROW EXECUTE FUNCTION update_timestamp_column();

-- CREATE TRIGGER comments_update_timestamp
-- BEFORE UPDATE ON comments
-- FOR EACH ROW EXECUTE FUNCTION update_timestamp_column();

-- CREATE TRIGGER moderation_queue_update_timestamp
-- BEFORE UPDATE ON moderation_queue
-- FOR EACH ROW EXECUTE FUNCTION update_timestamp_column();

-- CREATE TRIGGER reports_update_timestamp
-- BEFORE UPDATE ON reports
-- FOR EACH ROW EXECUTE FUNCTION update_timestamp_column();

-- CREATE TRIGGER scheduled_tasks_update_timestamp
-- BEFORE UPDATE ON scheduled_tasks
-- FOR EACH ROW EXECUTE FUNCTION update_timestamp_column();

-- CREATE TRIGGER contents_set_published_at
-- BEFORE UPDATE ON contents
-- FOR EACH ROW EXECUTE FUNCTION update_content_published_at();

-- -- 5. Índices
-- CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);
-- CREATE INDEX IF NOT EXISTS users_username_idx ON users(username);
-- CREATE INDEX IF NOT EXISTS works_author_id_idx ON works(author_id);
-- CREATE INDEX IF NOT EXISTS works_title_trgm_idx ON works USING gin (title gin_trgm_ops);
-- CREATE INDEX IF NOT EXISTS contents_work_id_idx ON contents(work_id);
-- CREATE INDEX IF NOT EXISTS content_images_content_id_idx ON content_images(content_id);
-- CREATE INDEX IF NOT EXISTS comments_content_id_idx ON comments(content_id);
-- CREATE INDEX IF NOT EXISTS likes_user_id_idx ON likes(user_id);
-- CREATE INDEX IF NOT EXISTS followers_followed_user_id_idx ON followers(followed_user_id);
-- CREATE INDEX IF NOT EXISTS notifications_user_id_idx ON notifications(user_id);
-- CREATE INDEX IF NOT EXISTS moderation_queue_status_idx ON moderation_queue(moderation_status);
-- CREATE INDEX IF NOT EXISTS reports_report_status_idx ON reports(report_status);
-- CREATE INDEX IF NOT EXISTS scheduled_tasks_status_idx ON scheduled_tasks(status);

-- -- 6. Views
-- CREATE VIEW users_safe AS
-- SELECT 
--     id,
--     username,
--     email,
--     role,
--     avatar,
--     verified,
--     status,
--     last_login_at,
--     created_at,
--     updated_at
-- FROM users;

-- -- 7. Dados Iniciais
-- INSERT INTO categories (name, slug, description) VALUES
-- ('Action', 'action', 'Stories with lots of action and adventure'),
-- ('Romance', 'romance', 'Stories focused on romantic relationships')
-- ON CONFLICT (name) DO NOTHING;

-- INSERT INTO genres (name, slug, description) VALUES
-- ('Shounen', 'shounen', 'Aimed at young male audience'),
-- ('Shoujo', 'shoujo', 'Aimed at young female audience')
-- ON CONFLICT (name) DO NOTHING;

-- COMMIT;