-- 启用向量插件（确保 PostgreSQL 安装了 pgvector）
CREATE EXTENSION IF NOT EXISTS vector;

-- 创建用户表
CREATE TABLE IF NOT EXISTS "user" (
                                      id SERIAL PRIMARY KEY,
                                      user_id VARCHAR(64) UNIQUE NOT NULL,         -- 分布式 ID，建议用雪花算法等生成
                                      username VARCHAR(256) UNIQUE,
                                      password TEXT NOT NULL,
                                      "like" TEXT,
                                      like_embedding VECTOR(128),                  -- 喜好向量，可根据接口长度调整维度
                                      create_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                      update_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- 索引（可选）
CREATE INDEX IF NOT EXISTS idx_user_user_id ON "user"(user_id);
CREATE INDEX IF NOT EXISTS idx_user_like_embedding ON "user" USING ivfflat (like_embedding vector_cosine_ops);
