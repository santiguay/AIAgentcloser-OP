# tienda/admin.py

from django.contrib import admin
from .models import Categoria, Producto, Orden, DetalleVenta


@admin.register(Categoria)
class CategoriaAdmin(admin.ModelAdmin):
    list_display = ('id', 'nombre', 'creada_en')
    search_fields = ('nombre',)
    ordering = ('nombre',)


@admin.register(Producto)
class ProductoAdmin(admin.ModelAdmin):
    list_display = ('id', 'nombre', 'categoria', 'precio', 'stock', 'creado_en')
    list_filter = ('categoria',)
    search_fields = ('nombre', 'categoria__nombre')
    ordering = ('nombre',)
    readonly_fields = ('creado_en',)


@admin.register(Orden)
class OrdenAdmin(admin.ModelAdmin):
    list_display = ('id', 'nombre_cliente', 'telefono', 'total', 'creada_en')
    search_fields = ('nombre_cliente', 'telefono', 'cedula')
    ordering = ('-creada_en',)
    readonly_fields = ('total', 'creada_en')
    list_filter = ('creada_en',)


@admin.register(DetalleVenta)
class DetalleVentaAdmin(admin.ModelAdmin):
    list_display = ('id', 'orden', 'producto', 'cantidad', 'subtotal')
    search_fields = ('orden__nombre_cliente', 'producto__nombre')
    ordering = ('orden',)
