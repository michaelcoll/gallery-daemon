CREATE TABLE photos
(
    hash          TEXT PRIMARY KEY,
    path          TEXT NOT NULL,
    date_time     TEXT,
    iso           INTEGER,
    exposure_time TEXT,
    x_dimension   INTEGER,
    y_dimension   INTEGER,
    model         TEXT,
    aperture      REAL
);
