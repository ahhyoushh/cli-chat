import time
import requests
import json
import os
from rich.console import Console
from rich.prompt import Prompt
from rich.table import Table

BASE_URL = "http://localhost:4000/api"  
CONFIG_FILE = "config.json"
SESSION_FILE = "session.json"

console = Console()

def load_config():
    if not os.path.exists(CONFIG_FILE):
        console.print(f"[bold red]Configuration file {CONFIG_FILE} not found.[/bold red]")
        raise FileNotFoundError(f"Configuration file {CONFIG_FILE} not found.")
    
    with open(CONFIG_FILE, "r") as file:
        return json.load(file)

def save_session(data):
    with open(SESSION_FILE, "w") as file:
        json.dump(data, file)

def load_session():
    if os.path.exists(SESSION_FILE):
        with open(SESSION_FILE, "r") as file:
            return json.load(file)
    return {}

def clear_session():
    if os.path.exists(SESSION_FILE):
        os.remove(SESSION_FILE)

def signup(username, password):
    url = f"{BASE_URL}/signup"
    payload = {"username": username, "password": password}
    response = requests.post(url, json=payload)
    message = response.text
    
    console.print(f"[bold green]{message}[/bold green]")

def login(username, password):
    url = f"{BASE_URL}/login"
    payload = {"username": username, "password": password}
    response = requests.post(url, json=payload)
    if response.status_code == 200:
        save_session({"username": username})
        
    else:
        console.print(response.text, style="bold red")
    return response.status_code == 200

def logout():
    clear_session()
    console.print("[bold green]Logged out successfully[/bold green]")

def is_logged_in():
    session = load_session()
    return "username" in session

def create_message(sender, receiver, message):
    if not is_logged_in():
        console.print("[bold red]You need to log in first.[/bold red]")
        return

    session = load_session()
    username = session["username"]

    url = f"{BASE_URL}/send"
    payload = {
        "sender": sender,
        "receiver": receiver,
        "message": message
    }
    response = requests.post(url, json=payload)
    if response.status_code == 200:
        msg = response.json()
        author = msg["sender"]
        receiver = msg["receiver"]
        message = msg["message"]
        console.print(f"[bold blue]{author} [/bold blue]said[green] '{message}'[/green] to [bold yellow]{receiver}![/bold yellow]")
    else:
        console.print(f'Check your args! {response.text}')
def get_all_messages():
    if not is_logged_in():
        console.print("[bold red]You need to log in first.[/bold red]")
        return

    session = load_session()
    username = session["username"]

    url = f"{BASE_URL}/getall"
    payload = {"username": username, "password": password}
    response = requests.post(url, json=payload)
    resS = response.text
    msgsList = json.loads(resS)
    for msg in msgsList:
        author = msg['sender']
        message = msg['message']
        console.print(f'[bold blue]{author}[/bold blue]: {message}')

def get_unread_messages():
    if not is_logged_in():
        console.print("[bold red]You need to log in first.[/bold red]")
        return

    session = load_session()
    username = session["username"]

    url = f"{BASE_URL}/getunreadmsgs"
    payload = {"username": username, "password": password}
    response = requests.post(url, json=payload)
    resS = response.text
    rest = resS[1:-2]
    try:
        msgsList = json.loads(rest)
        for msg in msgsList:
            author = msg['sender']
            message = msg['message']
            console.print(f'[bold blue]{author}[/bold blue]: {message}')
    except:
        console.print("[bold red] No new msgs![/bold red]")

if __name__ == "__main__":
    config = load_config()
    username = config["username"]
    password = config["password"]

    while True:
        console.print("[bold cyan]CLI Chat Application[/bold cyan]", style="bold blue")
        table = Table(show_header=True, header_style="bold magenta")
        table.add_column("Option", style="dim", width=12)
        table.add_column("Description")
        table.add_row("1", "Signup")
        table.add_row("2", "Login")
        table.add_row("3", "Logout")
        table.add_row("4", "Send Message")
        table.add_row("5", "Get All Messages")
        table.add_row("6", "Get Unread Messages")
        table.add_row("7", "Check if Logged In")
        table.add_row("8", "Exit")

        console.print(table)

        choice = Prompt.ask("[bold cyan]Enter your choice[/bold cyan]")

        if choice == "1":
            if username == "":
                print("Please set up your username and password in config.json!")
            signup(username, password)
            time.sleep(2)

        elif choice == "2":
            if login(username, password):
                console.print("[bold green]Login successful[/bold green]")
                time.sleep(2)
            else:
                console.print("[bold red]Login failed[/bold red]")
                time.sleep(3)

        elif choice == "3":
            logout()
            time.sleep(3)

        elif choice == "4":
            sender = config["username"]
            receiver = Prompt.ask("[bold cyan]Enter receiver username[/bold cyan]")
            message = Prompt.ask("[bold cyan]Enter message[/bold cyan]")
            create_message(sender, receiver, message)
            time.sleep(3)

        elif choice == "5":
            get_all_messages()
            break

        elif choice == "6":
            get_unread_messages()
            time.sleep(1)

        elif choice == "7":
            if is_logged_in():
                console.print("[bold green]You are logged in.[/bold green]")
                time.sleep(1)
            else:
                console.print("[bold red]You are not logged in.[/bold red]")
                time.sleep(1)

        elif choice == "8":
            break

        else:
            console.print("[bold red]Invalid choice[/bold red]")
            time.sleep(1)
