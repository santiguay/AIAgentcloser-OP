const { createBot, createProvider, createFlow, addKeyword } = require('@bot-whatsapp/bot');
const QRPortalWeb = require('@bot-whatsapp/portal');
const BaileysProvider = require('@bot-whatsapp/provider/baileys');
const MockAdapter = require('@bot-whatsapp/database/mock');
const { GoogleGenerativeAI } = require('@google/generative-ai');
const axios = require('axios');

const API_BASE_URL = 'http://localhost:8082';

require('dotenv').config();


const { GoogleGenerativeAI } = require('@google/generative-ai');


const genAI = new GoogleGenerativeAI(process.env.GOOGLE_API_KEY);
const conversationHistory = new Map(); 
// Función para obtener o crear un chat
async function getOrCreateChat(phoneNumber) {
    if (!conversationHistory.has(phoneNumber)) {
        const model = genAI.getGenerativeModel({ model: 'gemini-pro' });
        const chat = model.startChat({
            history: [],
            generationConfig: {
                maxOutputTokens: 2000,
                temperature: 0.9,
                topP: 0.8,
                topK: 40,
            },
        });

        conversationHistory.set(phoneNumber, {
            chat,
            lastAccess: Date.now(),
        });
    }
    return conversationHistory.get(phoneNumber);
}

// Función para enviar un mensaje a Gemini y manejar respuestas
async function getChatResponse(phoneNumber, message) {
    try {
        prompt1 = `Eres un vendedor de ropa de la empresa "Virtual Dress". Tu trabajo es vender productos y responder preguntas relacionadas con el negocio y cerrar las ventas. 
Si detectas interés en categorías, responde solo el texto y unicamente "mostrar_categorias". 
Si detectas interés en productos, responde solo el texto y unicamente"listar_productos". 
Si tienes todos los datos para crear un pedido, responde solo el texto y unicamente"crear_orden" junto con el JSON del pedido(con este formato es el JSON de la orden: {"order":{"nombre_cliente": "Juan Pérez","domicilio": "Calle Falsa 123","cedula": "1234567890","telefono": "555-1234"},"detalle_venta":[{"producto_id":1,"cantidad":2},{"producto_id":2,"cantidad":3}]}).
Si no sabes la respuesta, di que no puedes responder preguntas fuera de este tema.
Sé marketero, haz recomendaciones, da información como "Número de contacto: 555-123-456", y mantén la conversación enfocada en ventas.
Y el mensaje es: ${message}`
        const conversation = await getOrCreateChat(phoneNumber);
        const result = await conversation.chat.sendMessage(prompt1);
        const response = result.response;

  
        if (response && Array.isArray(response.candidates) && response.candidates[0]?.content) {
            const textResponse = response.candidates[0].content; // Extraer contenido principal
            conversation.lastAccess = Date.now(); // Actualizar acceso reciente
            return typeof textResponse === 'string' ? textResponse.trim() : textResponse["parts"][0]["text"];
        }

        console.error('Respuesta inesperada de Gemini:', response);
        return 'Lo siento, hubo un problema al interpretar tu mensaje. ¿Podrías intentar de nuevo?';
    } catch (error) {
        console.error('Error al obtener respuesta:', error);


        if (error.message.includes('too long')) {
            conversationHistory.delete(phoneNumber);
            return 'Lo siento, tuve que reiniciar nuestra conversación porque se hizo muy larga. ¿Podrías repetir tu última pregunta?';
        }

        return 'Lo siento, hubo un error al procesar tu mensaje. ¿Podrías intentar de nuevo?';
    }
}



async function apiRequest(endpoint, method = 'GET', data = null) {
    try {
        const response = await axios({
            url: `${API_BASE_URL}${endpoint}`,
            method,
            data,
        });
        return response.data;
    } catch (error) {
        console.error(`Error en la petición ${method} ${endpoint}:`, error);
        return null;
    }
}

// Flujo principal
const flowPrincipal = addKeyword([''])
    .addAction(async (ctx, { flowDynamic }) => {
        if (!ctx.body) return;

        try {
            const phoneNumber = ctx.from;
            console.log(`Mensaje recibido de ${phoneNumber}: ${ctx.body}`); // Debugging

            const aiResponse = await getChatResponse(phoneNumber, ctx.body);
            console.log(`Respuesta AI para ${phoneNumber}: ${aiResponse}`); // Debugging

            if (aiResponse === 'mostrar_categorias') {
                const categorias = await apiRequest('/categorias');
                if (categorias) {
                    const lista = categorias.map(c => `- ${c.nombre}`).join('\n');
                    return await flowDynamic(`Estas son las categorías disponibles:\n${lista}`);
                }
                return await flowDynamic('No se pudieron obtener las categorías en este momento.');
            }

            if (aiResponse === 'listar_productos') {
                const productos = await apiRequest('/productos');
                if (productos) {
                    const lista = productos.map(p => `- ${p.nombre}: $${p.precio} (Stock: ${p.stock})`).join('\n');
                    return await flowDynamic(`Estos son los productos disponibles:\n${lista}`);
                }
                return await flowDynamic('No se pudieron obtener los productos en este momento.');
            }

            if (aiResponse.startsWith('crear_orden')) {
                try {
                    const datosOrden = JSON.parse(aiResponse.replace('crear_orden', '').trim());
                    const nuevaOrden = await apiRequest('/ordenes', 'POST', datosOrden);
                    if (nuevaOrden) {
                        return await flowDynamic(`¡Pedido creado con éxito! Número de pedido: ${nuevaOrden.id}.`);
                    }
                    return await flowDynamic('No se pudo crear el pedido en este momento.');
                } catch (error) {
                    console.error('Error procesando la orden:', error);
                    return await flowDynamic('Hubo un error procesando tu pedido. Por favor, intenta nuevamente.');
                }
            }

     
            return await flowDynamic(aiResponse);
        } catch (error) {
            console.error('Error en el flujo principal:', error);
            return await flowDynamic('Ocurrió un error. Por favor, intenta nuevamente.');
        }
    });

// Inicialización del bot
const main = async () => {
    const adapterDB = new MockAdapter();
    const adapterFlow = createFlow([flowPrincipal]);
    const adapterProvider = createProvider(BaileysProvider);

    createBot({
        flow: adapterFlow, 
        provider: adapterProvider,
        database: adapterDB,
    });

    QRPortalWeb(); 
};

main();
