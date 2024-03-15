import subprocess

long_urls = [
    "https://cdn.ttgtmedia.com/rms/onlineimages/networking-cdn_configuration.png",
    "https://example.com",
    "https://images.pexels.com/photos/674010/pexels-photo-674010.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1",
    "https://developers.cloudflare.com/assets/private-ips-diagram_hua0d86fa4f03f384bde81e46921dab3a3_510169_2530x1144_resize_q75_box_3-8ce36726.png"

]


def generate_short_url(long_url: str) -> str:
    cmd = f"""curl --location 'http://localhost:8080/api/new' \
    --header 'Content-Type: application/json' \
    --data '{{
        "full_url": "{long_url}"
    }}'"""
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    return result.stdout


def get_data_from_url(url: str) -> bytes:
    cmd = f"""curl --location '{url}'"""
    result = subprocess.run(cmd, shell=True, capture_output=True)
    return result.stdout


long_to_short = {long_url: "http://localhost:8080/" + generate_short_url(long_url) for long_url in long_urls}


for long_url, short_url in long_to_short.items():
    if get_data_from_url(long_url) == get_data_from_url(short_url):
        print(f"Test passed for {long_url}")
    else:
        print(f"Test failed for {long_url}")

