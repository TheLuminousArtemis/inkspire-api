ALTER TABLE comments
ADD COLUMN parent_id bigint DEFAULT NULL,
ADD COLUMN deleted boolean DEFAULT FALSE,
ADD CONSTRAINT fk_comment_parent
FOREIGN KEY (parent_id) REFERENCES public.comments(id) ON DELETE CASCADE;