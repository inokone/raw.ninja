import inject

from rawninja.dependency_register import register_dependencies
from rawninja.handler import ExpiryMarker

register_dependencies()

handler = inject.instance(ExpiryMarker)
