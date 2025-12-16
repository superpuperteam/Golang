CREATE TABLE IF NOT EXISTS movies (
                                      id    SERIAL PRIMARY KEY,
                                      title TEXT NOT NULL,
                                      year  INT  NOT NULL
);

CREATE TABLE IF NOT EXISTS actors (
                                      id       SERIAL PRIMARY KEY,
                                      movie_id INT NOT NULL REFERENCES movies(id),
    name     TEXT NOT NULL
    );

TRUNCATE actors, movies RESTART IDENTITY;

INSERT INTO movies (title, year) VALUES
                                     ('Inception', 2010),
                                     ('Interstellar', 2014),
                                     ('Mad Max: Fury Road', 2015),
                                     ('The Dark Knight', 2008);

INSERT INTO actors (movie_id, name) VALUES
                                        (1,'Leonardo DiCaprio'),
                                        (1,'Joseph Gordon-Levitt'),
                                        (1,'Elliot Page'),
                                        (1,'Tom Hardy'),
                                        (1,'Ken Watanabe'),
                                        (2,'Matthew McConaughey'),
                                        (2,'Anne Hathaway'),
                                        (2,'Jessica Chastain'),
                                        (2,'Michael Caine'),
                                        (3,'Tom Hardy'),
                                        (3,'Charlize Theron');
