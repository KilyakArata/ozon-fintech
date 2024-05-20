INSERT INTO posts (author, body, open) VALUES ('bob', 'interesting post', TRUE);
INSERT INTO posts (author, body, open) VALUES ('ann', 'beautiful post', TRUE);
INSERT INTO posts (author, body, open) VALUES ('jon', 'fantastic post', TRUE);

INSERT INTO comments (author, body, post_id) VALUES ('ann', 'wow', 1);
INSERT INTO comments (author, body, post_id) VALUES ('bob', 'heh', 3);
INSERT INTO comments (author, body, post_id) VALUES ('jon', 'yeah', 2);

INSERT INTO comments (author, body, post_id, parent_comment_id) VALUES ('sis', 'rer', 1, 1);
INSERT INTO comments (author, body, post_id, parent_comment_id) VALUES ('mim', 'geg', 1, 1);
INSERT INTO comments (author, body, post_id, parent_comment_id) VALUES ('nin', 'gqw', 1, 5);