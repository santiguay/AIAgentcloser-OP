```markdown
# AI Agent for Clothing Store

This project integrates an **Artificial Intelligence agent** designed to manage orders for a clothing store. The project includes:

- **Django Admin Panel**: Located in the `clothing_store` folder, where orders are received and managed.
- **Go API**: Developed so that the agent consumes a single API endpoint, preventing SQL injection.
- **Chatbot using Baileys**: Implemented to interact via WhatsApp. *(Credit to [Baileys](https://github.com/adiwajshing/Baileys) for the tool!)*

---

## Project Structure

- **`clothing_store` Folder**:  
  This folder serves as the administrator area where orders are received and managed.  
  **Subfolder `storeSettings`**:  
  Within this subfolder, create the `settings.py` file containing the Django configuration.

- **Go API**:  
  A dedicated API that ensures secure data consumption and prevents SQL injection.

- **Chatbot**:  
  Developed using Baileys. Install dependencies with npm and then run the application.

---

## Django Admin Panel Configuration

Within the `storeSettings` subfolder, create the `settings.py` file and install the necessary requirements. Below is an example configuration:

```python
INSTALLED_APPS = [
    'unfold',
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
    'adminApp',
]

# Admin template settings
UNFOLD = {
    "site_title": "Clothing Store",
    "site_header": "Administration Panel",
    "welcome_sign": "Welcome to the Clothing Store Admin Panel",
    # Model used for the global search bar
}

MIDDLEWARE = [
    'django.middleware.security.SecurityMiddleware',
    'django.contrib.sessions.middleware.SessionMiddleware',
    'django.middleware.common.CommonMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'django.middleware.clickjacking.XFrameOptionsMiddleware',
]

ROOT_URLCONF = 'storeSettings.urls'

TEMPLATES = [
    {
        'BACKEND': 'django.template.backends.django.DjangoTemplates',
        'DIRS': [],
        'APP_DIRS': True,
        'OPTIONS': {
            'context_processors': [
                'django.template.context_processors.debug',
                'django.template.context_processors.request',
                'django.contrib.auth.context_processors.auth',
                'django.contrib.messages.context_processors.messages',
            ],
        },
    },
]

WSGI_APPLICATION = 'storeSettings.wsgi.application'
```

---

## Installation and Execution

### 1. Django Setup

1. **Clone the Repository**  
   Clone the project to your local machine.

2. **Create and Activate a Virtual Environment**  
   Make sure you have Python installed and create a virtual environment:
   ```bash
   python -m venv env
   source env/bin/activate  # On Linux/Mac
   env\Scripts\activate     # On Windows
   ```

3. **Install Dependencies**  
   Install the required packages:
   ```bash
   pip install -r requirements.txt
   ```

4. **Run Migrations and Start the Server**  
   Navigate to the project's root or the `clothing_store` folder and execute:
   ```bash
   python manage.py migrate
   python manage.py runserver
   ```

### 2. Running the Go API

The AI agent consumes a dedicated API built in Go, ensuring the use of a single endpoint to prevent SQL injections. Please follow the specific instructions for the Go service to:

- Set up the environment.
- Run the API.

> **Note:** Refer to the Go API documentation for details about endpoints and necessary configurations.

### 3. Setting Up and Running the Chatbot with Baileys

The chatbot is built using [Baileys](https://github.com/adiwajshing/Baileys). To run it:

1. Navigate to the chatbot folder.
2. Install the dependencies using npm:
   ```bash
   npm install
   ```
3. Start the chatbot:
   ```bash
   node app.js
   ```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```
