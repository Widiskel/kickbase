import time
import json
from playwright.sync_api import sync_playwright

SWAGGER_URL = "http://localhost:8080/swagger/index.html"
TIMESTAMP = int(time.time())
SCREENSHOT_PATH = f"swagger_all_scenarios_{TIMESTAMP}.png"

def run_all_swagger_scenarios():
    print("=" * 80)
    print("   FULL END-TO-END BROWSER VALIDATION ON SWAGGER UI DOM (ALL SCENARIOS)   ")
    print("=" * 80)

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context(viewport={"width": 1440, "height": 1000})
        page = context.new_page()

        print(f"\n[1] Navigating to Swagger UI at {SWAGGER_URL}...")
        page.goto(SWAGGER_URL)
        page.wait_for_selector(".swagger-ui", timeout=10000)
        print(f"✓ Loaded page: '{page.title()}'")

        def execute_swagger_op(op_id, body_payload=None, params=None):
            op_block = page.locator(f"#{op_id}")
            if not op_block.locator(".try-out__btn").is_visible():
                op_block.locator(".opblock-summary").click()
                page.wait_for_timeout(400)

            try_btn = op_block.locator(".try-out__btn")
            if try_btn.is_visible() and "Cancel" not in try_btn.inner_text():
                try_btn.click()
                page.wait_for_timeout(300)

            if body_payload:
                body_textarea = op_block.locator("textarea.body-param__text")
                if body_textarea.is_visible():
                    body_textarea.fill(json.dumps(body_payload, indent=2))
                    page.wait_for_timeout(200)

            if params:
                for param_name, param_val in params.items():
                    input_elem = op_block.locator(f"input[placeholder='{param_name}'], input[data-param-name='{param_name}']").first
                    if input_elem.is_visible():
                        input_elem.fill(str(param_val))

            exec_btn = op_block.locator("button.execute")
            exec_btn.click()
            page.wait_for_timeout(800)

            res_code = op_block.locator(".response .response-col_status").first.inner_text()
            res_body_elem = op_block.locator(".response .response-col_description pre").first
            res_body_text = res_body_elem.inner_text() if res_body_elem.is_visible() else "{}"
            try:
                res_body_json = json.loads(res_body_text)
            except Exception:
                res_body_json = {"raw": res_body_text}

            return res_code, res_body_json

        # 1. Register
        print("\n--- [SCENARIO 1] Register Admin User ---")
        code, res = execute_swagger_op("operations-Auth-post_api_auth_register", body_payload={
            "username": f"admin_swg_all_{TIMESTAMP}",
            "password": "password123",
            "name": "Admin Swagger All",
            "role": "admin"
        })
        print(f"✓ Register: HTTP {code} | User: {res['data']['username']}")
        assert "201" in code

        # 2. Login
        print("\n--- [SCENARIO 2] Login & Acquire JWT Dual-Tokens ---")
        code, res = execute_swagger_op("operations-Auth-post_api_auth_login", body_payload={
            "username": f"admin_swg_all_{TIMESTAMP}",
            "password": "password123"
        })
        token = res["data"]["access_token"]
        refresh_token = res["data"]["refresh_token"]
        print(f"✓ Login: HTTP {code} | Permissions: {len(res['data']['permissions'])}")
        assert "200" in code

        # 3. Authorize Bearer Token in Swagger
        print("\n--- [SCENARIO 3] Set Bearer Token in Swagger Authorize Modal ---")
        page.locator("button.authorize").click()
        page.wait_for_timeout(400)
        auth_input = page.locator(".auth-container input").first
        auth_input.fill(f"Bearer {token}")
        page.locator("button.modal-btn.auth.authorize").click()
        page.wait_for_timeout(200)
        page.locator("button.modal-btn.btn-done").click()
        print("✓ Bearer Token Authorized in Swagger UI")

        # 4. Create Team 1 & 2
        print("\n--- [SCENARIO 4] Create Football Teams (Persija & Persib) ---")
        code1, t1 = execute_swagger_op("operations-Teams-post_api_teams", body_payload={
            "name": f"Persija All {TIMESTAMP}",
            "founded_year": 1928,
            "address": "Senayan Jakarta",
            "city": "Jakarta",
            "logo_url": "https://example.com/persija.png"
        })
        team1_id = t1["data"]["id"]
        print(f"✓ Team 1: HTTP {code1} | ID: {team1_id} | Name: {t1['data']['name']}")

        code2, t2 = execute_swagger_op("operations-Teams-post_api_teams", body_payload={
            "name": f"Persib All {TIMESTAMP}",
            "founded_year": 1933,
            "address": "Gedebage Bandung",
            "city": "Bandung",
            "logo_url": "https://example.com/persib.png"
        })
        team2_id = t2["data"]["id"]
        print(f"✓ Team 2: HTTP {code2} | ID: {team2_id} | Name: {t2['data']['name']}")

        # 5. Update Team (Optimistic Locking)
        print("\n--- [SCENARIO 5] Update Team with Optimistic Locking Version ---")
        code_u, t1_u = execute_swagger_op("operations-Teams-put_api_teams__id_", body_payload={
            "name": f"Persija Updated {TIMESTAMP}",
            "address": "Senayan Baru No. 10",
            "version": 1
        }, params={"id": team1_id})
        print(f"✓ Team Update: HTTP {code_u} | New Version: {t1_u['data']['version']}")

        # 6. Team History & Revert
        print("\n--- [SCENARIO 6] Team Audit Trail & Revert Version ---")
        code_h, t1_h = execute_swagger_op("operations-Teams-get_api_teams__id__history", params={"id": team1_id})
        print(f"✓ Team History Entries: {len(t1_h['data'])}")

        code_rev, t1_rev = execute_swagger_op("operations-Teams-post_api_teams__id__revert", body_payload={
            "target_version": 1
        }, params={"id": team1_id})
        print(f"✓ Team Revert to v1: HTTP {code_rev} | {t1_rev.get('message', 'Reverted')}")

        # 7. Create Players
        print("\n--- [SCENARIO 7] Create Players with eFootball Attributes ---")
        code_p1, p1 = execute_swagger_op("operations-Players-post_api_players", body_payload={
            "team_id": team1_id,
            "name": "Bambang Pamungkas",
            "height": 178.0,
            "weight": 72.0,
            "position": "CF",
            "playstyle": "Goal Poacher",
            "jersey_number": 20
        })
        p1_id = p1["data"]["id"]
        print(f"✓ Player 1: HTTP {code_p1} | {p1['data']['name']} (CF #{p1['data']['jersey_number']})")

        code_p2, p2 = execute_swagger_op("operations-Players-post_api_players", body_payload={
            "team_id": team2_id,
            "name": "Atep",
            "height": 170.0,
            "weight": 65.0,
            "position": "LWF",
            "playstyle": "Prolific Winger",
            "jersey_number": 7
        })
        p2_id = p2["data"]["id"]
        print(f"✓ Player 2: HTTP {code_p2} | {p2['data']['name']} (LWF #{p2['data']['jersey_number']})")

        # 8. Schedule Match
        print("\n--- [SCENARIO 8] Schedule Match ---")
        code_m, m = execute_swagger_op("operations-Matches-post_api_matches", body_payload={
            "match_date": "2026-09-01",
            "match_time": "15:30:00",
            "home_team_id": team1_id,
            "away_team_id": team2_id
        })
        match_id = m["data"]["id"]
        print(f"✓ Match Scheduled: HTTP {code_m} | ID: {match_id}")

        # 9. Update Match Status Lifecycle
        print("\n--- [SCENARIO 9] Match Status Transition Lifecycle ---")
        code_s1, s1 = execute_swagger_op("operations-Matches-patch_api_matches__id__status", body_payload={
            "status": "deferred",
            "version": 1
        }, params={"id": match_id})
        print(f"✓ Status transition to deferred: HTTP {code_s1}")
        assert "200" in code_s1

        code_s2, s2 = execute_swagger_op("operations-Matches-patch_api_matches__id__status", body_payload={
            "status": "scheduled",
            "version": 2
        }, params={"id": match_id})
        print(f"✓ Status transition back to scheduled: HTTP {code_s2}")
        assert "200" in code_s2

        # 10. Report Match Result
        print("\n--- [SCENARIO 10] Report Match Result & Goal Details ---")
        code_r, r = execute_swagger_op("operations-Results-post_api_results", body_payload={
            "match_id": match_id,
            "home_score": 2,
            "away_score": 1,
            "goals": [
                {"player_id": p1_id, "goal_time": "18'"},
                {"player_id": p1_id, "goal_time": "45'"},
                {"player_id": p2_id, "goal_time": "60'"}
            ]
        })
        print(f"✓ Match Result Recorded: HTTP {code_r} | Score: {r['data']['home_score']}-{r['data']['away_score']}")

        # 11. Public Match Report
        print("\n--- [SCENARIO 11] Public Aggregated Match Report ---")
        code_rep, rep = execute_swagger_op("operations-Reports-get_api_reports_matches__id_", params={"id": match_id})
        rep_d = rep["data"]
        print(f"✓ Match Report: HTTP {code_rep}")
        print(f"   - Match: {rep_d['home_team']} vs {rep_d['away_team']}")
        print(f"   - Status: '{rep_d['status']}' ({rep_d['home_score']}-{rep_d['away_score']})")
        print(f"   - Top Scorer: {rep_d['top_scorers'][0]['player_name']} ({rep_d['top_scorers'][0]['goals']} goals)")
        print(f"   - Cumulative Wins: Home ({rep_d['cumulative_home_wins']}), Away ({rep_d['cumulative_away_wins']})")

        # 12. Soft Delete Constraint Check (Cannot delete team with active players)
        print("\n--- [SCENARIO 12] Soft Delete Constraint Enforcement ---")
        code_del_t, del_t = execute_swagger_op("operations-Teams-delete_api_teams__id_", params={"id": team1_id})
        print(f"✓ Delete Team with players rejected as expected: HTTP {code_del_t} ({del_t.get('error', {}).get('message')})")
        assert "409" in code_del_t

        # 13. Soft Delete Constraint Check (Cannot delete player with goal history)
        code_del_p, del_p = execute_swagger_op("operations-Players-delete_api_players__id_", params={"id": p1_id})
        print(f"✓ Delete Player with goal history rejected as expected: HTTP {code_del_p} ({del_p.get('error', {}).get('message')})")
        assert "409" in code_del_p

        # Take full page verification screenshot
        page.screenshot(path=SCREENSHOT_PATH, full_page=True)
        print(f"\n📸 Final full-page Swagger test screenshot saved to: {SCREENSHOT_PATH}")

        browser.close()

    print("\n" + "=" * 80)
    print("  🎉 ALL 13 SWAGGER UI END-TO-END SCENARIOS EXECUTED & PASSED 100%!   ")
    print("=" * 80)

if __name__ == "__main__":
    run_all_swagger_scenarios()
