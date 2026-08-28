-- 1. Seed Default Admin, Staff, and Viewer (Password: password123)
-- bcrypt hash for 'password123': $2a$10$iIuImsC8Z7m/E0R.fXz5gO4tF3dG3oW3G5X8K1g/vW8b6F3qY4j8e
INSERT INTO users (id, username, password, name, role, created_at, updated_at)
VALUES 
  ('11111111-1111-1111-1111-111111111111', 'admin', '$2a$10$p4xL4M0O5nS1f9uU7T.3fO0F0BqO0Y5K1G/vW8b6F3qY4j8eXXXXX', 'System Administrator', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('22222222-2222-2222-2222-222222222222', 'staff', '$2a$10$p4xL4M0O5nS1f9uU7T.3fO0F0BqO0Y5K1G/vW8b6F3qY4j8eXXXXX', 'Club Operator Staff', 'staff', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('33333333-3333-3333-3333-333333333333', 'viewer', '$2a$10$p4xL4M0O5nS1f9uU7T.3fO0F0BqO0Y5K1G/vW8b6F3qY4j8eXXXXX', 'Public Viewer', 'viewer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- 2. Seed Sample Teams
INSERT INTO teams (id, name, founded_year, address, city, logo_url, version, created_at, updated_at)
VALUES
  ('a1111111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Persija Jakarta', 1928, 'Jl. Pintu Satu Senayan', 'Jakarta', 'https://upload.wikimedia.org/wikipedia/id/5/5e/Logo_Persija.png', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('b2222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Persib Bandung', 1933, 'Jl. Sulanjana No. 17', 'Bandung', 'https://upload.wikimedia.org/wikipedia/id/8/80/Persib_Bandung_logo.svg', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('c3333333-cccc-cccc-cccc-cccccccccccc', 'Arema FC', 1987, 'Jl. Mayjend Panjaitan No. 42', 'Malang', 'https://upload.wikimedia.org/wikipedia/id/2/2b/Arema_FC_logo.png', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('d4444444-dddd-dddd-dddd-dddddddddddd', 'Persebaya Surabaya', 1927, 'Jl. Embong Malang No. 88', 'Surabaya', 'https://upload.wikimedia.org/wikipedia/id/0/07/Persebaya_Surabaya_logo.svg', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- 3. Seed Sample Players
INSERT INTO players (id, team_id, name, height, weight, position, playstyle, jersey_number, version, created_at, updated_at)
VALUES
  -- Persija
  ('p1111111-1111-1111-1111-111111111111', 'a1111111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Bambang Pamungkas', 178, 72, 'CF', 'Goal Poacher', 20, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('p1111111-1111-1111-1111-222222222222', 'a1111111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Andritany Ardhiyasa', 180, 75, 'GK', 'Defensive Goalkeeper', 26, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  -- Persib
  ('p2222222-2222-2222-2222-111111111111', 'b2222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Atep', 170, 65, 'LWF', 'Prolific Winger', 7, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  ('p2222222-2222-2222-2222-222222222222', 'b2222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'I Made Wirawan', 181, 76, 'GK', 'Defensive Goalkeeper', 78, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;

-- 4. Seed Initial Match
INSERT INTO matches (id, match_date, match_time, home_team_id, away_team_id, status, version, created_at, updated_at)
VALUES
  ('m1111111-1111-1111-1111-111111111111', '2026-09-01', '15:30:00', 'a1111111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'b2222222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'scheduled', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;
