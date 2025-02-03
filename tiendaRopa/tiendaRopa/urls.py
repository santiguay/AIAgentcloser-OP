"""
URL configuration for tiendaRopa project.

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/5.1/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path


from django.urls import path
from adminApp.views import PendingOrdersView, OrderDetailView, CompleteOrderView
from django.views.generic import TemplateView

urlpatterns = [
    path('', TemplateView.as_view(template_name='Ordenes.html'), name='pending_orders'),
    path('admin/', admin.site.urls),
    path('api/pending-orders/', PendingOrdersView.as_view(), name='api-pending-orders'),
    path('api/orders/<int:pk>/', OrderDetailView.as_view(), name='api-order-detail'),
    path('api/complete-order/<int:pk>/', CompleteOrderView.as_view(), name='api-complete-order'),
]