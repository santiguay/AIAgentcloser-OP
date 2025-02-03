from django.shortcuts import render

from django.views import View
from django.http import JsonResponse
from django.shortcuts import get_object_or_404
from .models import Orden
import json
from django.utils.dateparse import parse_datetime
from django.utils.timezone import is_aware, make_aware
from django.conf import settings
from django.utils.dateparse import parse_datetime
from datetime import datetime
class PendingOrdersView(View):
    def get(self, request):
        orders = Orden.objects.filter(completa=False).order_by('-creada_en')
        data = []
        for order in orders:
            created_at = order.creada_en
            if isinstance(created_at, str):
                # If it's a string, try to parse it into a datetime object
                created_at = parse_datetime(created_at)
                if created_at is not None and settings.USE_TZ and not is_aware(created_at):
                    # Make the datetime timezone-aware
                    created_at = make_aware(created_at)
            
            # If parsing failed or it's already a datetime, format it
            if isinstance(created_at, datetime):
                formatted_date = created_at.strftime('%Y-%m-%d %H:%M:%S')
            else:
                # If all else fails, use the original string
                formatted_date = str(order.creada_en)

            data.append({
                'id': order.id,
                'nombre_cliente': order.nombre_cliente,
                'total': str(order.total),
                'creada_en': formatted_date
            })
        return JsonResponse(data, safe=False)

class OrderDetailView(View):
    def get(self, request, pk):
        order = get_object_or_404(Orden, pk=pk)
        data = {
            'id': order.id,
            'nombre_cliente': order.nombre_cliente,
            'domicilio': order.domicilio,
            'cedula': order.cedula,
            'telefono': order.telefono,
            'total': str(order.total),
            'completa': order.completa,
            'creada_en': order.creada_en,
            'detalles': [{
                'producto': detalle.producto.nombre,
                'cantidad': detalle.cantidad,
                'subtotal': str(detalle.subtotal)
            } for detalle in order.detalles.all()]
        }
        return JsonResponse(data)

from django.views import View
from django.http import JsonResponse
from django.shortcuts import get_object_or_404
from django.db import transaction
from .models import Orden, Producto
from django.core.exceptions import ValidationError
from django.views.decorators.csrf import csrf_exempt
from django.utils.decorators import method_decorator

@method_decorator(csrf_exempt, name='dispatch')
class CompleteOrderView(View):
    @transaction.atomic
    def post(self, request, pk):
        try:
            order = get_object_or_404(Orden, pk=pk)
            
            if order.completa:
                return JsonResponse({'message': 'La orden ya está completada'}, status=400)

            for detalle in order.detalles.all():
                producto = detalle.producto
                if producto.stock < detalle.cantidad:
                    raise ValidationError(f"Stock insuficiente para {producto.nombre}")
                
                producto.stock -= detalle.cantidad
                producto.save()

            order.completa = True
            order.save()

            return JsonResponse({'message': 'Orden completada exitosamente y stock actualizado'})
        
        except ValidationError as e:
            return JsonResponse({'error': str(e)}, status=400)
        except Exception as e:
            return JsonResponse({'error': 'Error al completar la orden'}, status=500)
    
    
  

