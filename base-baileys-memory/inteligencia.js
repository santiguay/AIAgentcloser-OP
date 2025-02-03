// Cargar variables del archivo .env
require('dotenv').config();

// Importar la librería de Google Generative AI
const { GoogleGenerativeAI } = require('@google/generative-ai');

// Configurar Gemini usando la clave API de las variables de entorno
const genAI = new GoogleGenerativeAI(process.env.GOOGLE_API_KEY);


// Contexto inicial para Gemini
const CONTEXTO_INICIAL = `
Eres un vendedor de ropa de la empresa "Virtual Dress". Tu trabajo es vender productos y responder preguntas relacionadas con las categorías y productos. 
Si detectas interés en categorías, responde solo el texto y unicamente "mostrar_categorias". 
Si detectas interés en productos, responde solo el texto y unicamente"listar_productos". 
Si tienes todos los datos para crear un pedido, responde solo el texto y unicamente"crear_orden" junto con el JSON del pedido.
Si no sabes la respuesta, di que no puedes responder preguntas fuera de este tema.
Sé marketero, haz recomendaciones, da información como "Número de contacto: 555-123-456", y mantén la conversación enfocada en ventas.
`;

// Mapa para almacenar conversaciones activas
const conversationHistory = new Map();

// Obtener o crear un chat para un cliente
async function getOrCreateChat(phoneNumber) {
    if (!conversationHistory.has(phoneNumber)) {
        const model = genAI.getGenerativeModel({ model: 'gemini-pro' });
        const chat = model.startChat({
            history: [{ role: 'vendedor de ropa de la empresa "Virtual Dress"', content: CONTEXTO_INICIAL }],
            generationConfig: {
                maxOutputTokens: 2000,
                temperature: 0.7,
                topP: 0.9,
                topK: 40,
            },
        });
        conversationHistory.set(phoneNumber, chat);
    }
    return conversationHistory.get(phoneNumber);
}

// Interpretar el mensaje
async function interpretarMensaje(phoneNumber, mensaje) {
    const chat = await getOrCreateChat(phoneNumber);

    // Enviar el mensaje del usuario al chat
    chat.history.push({ role: 'cliente', content: mensaje });
    const response = await chat.sendMessage(mensaje);
    const aiResponse = await response.response;

    // Devolver la respuesta interpretada
    return aiResponse.text.trim();
}

module.exports = { interpretarMensaje };
