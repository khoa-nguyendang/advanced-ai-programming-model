# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: app/main/grpc/faiss.proto

from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='app/main/grpc/faiss.proto',
  package='faiss',
  syntax='proto3',
  serialized_options=b'Z#services/user/protos/faiss/v1;faiss',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x19\x61pp/main/grpc/faiss.proto\x12\x05\x66\x61iss\")\n\x06Vector\x12\x11\n\tvector_id\x18\x01 \x01(\t\x12\x0c\n\x04\x64\x61ta\x18\x02 \x03(\x01\"K\n\x17VectorEnrollmentRequest\x12\x10\n\x08user_uid\x18\x01 \x01(\t\x12\x1e\n\x07vectors\x18\x02 \x03(\x0b\x32\r.faiss.Vector\"<\n\x15VectorDeletionRequest\x12\x11\n\timage_ids\x18\x01 \x03(\t\x12\x10\n\x08user_uid\x18\x02 \x01(\t\"L\n\x18VectorEnrollmentResponse\x12\x1f\n\x04\x63ode\x18\x01 \x01(\x0e\x32\x11.faiss.StatusCode\x12\x0f\n\x07message\x18\x02 \x01(\t\"5\n\x13VectorSearchRequest\x12\x1e\n\x07vectors\x18\x01 \x03(\x0b\x32\r.faiss.Vector\"{\n\x14VectorSearchResponse\x12\x1f\n\x04\x63ode\x18\x01 \x01(\x0e\x32\x11.faiss.StatusCode\x12\x0f\n\x07message\x18\x02 \x01(\t\x12\x10\n\x08user_uid\x18\x03 \x01(\t\x12\x10\n\x08image_id\x18\x04 \x01(\t\x12\r\n\x05score\x18\x05 \x01(\x01\"J\n\x16VectorDeletionResponse\x12\x1f\n\x04\x63ode\x18\x01 \x01(\x0e\x32\x11.faiss.StatusCode\x12\x0f\n\x07message\x18\x02 \x01(\t*\x8b\x01\n\nStatusCode\x12\x10\n\x0c\x45NGINE_ERROR\x10\x00\x12\x0b\n\x07SUCCESS\x10\x01\x12\r\n\tBAD_IMAGE\x10\x02\x12\x0e\n\nBAD_VECTOR\x10\x03\x12%\n!IMAGE_REGISTERED_FOR_ANOTHER_USER\x10\x04\x12\t\n\x05\x46OUND\x10\x05\x12\r\n\tNOT_FOUND\x10\x06\x32\xe0\x01\n\tVectorAPI\x12I\n\x06\x45nroll\x12\x1e.faiss.VectorEnrollmentRequest\x1a\x1f.faiss.VectorEnrollmentResponse\x12\x41\n\x06Search\x12\x1a.faiss.VectorSearchRequest\x1a\x1b.faiss.VectorSearchResponse\x12\x45\n\x06\x44\x65lete\x12\x1c.faiss.VectorDeletionRequest\x1a\x1d.faiss.VectorDeletionResponseB%Z#services/user/protos/faiss/v1;faissb\x06proto3'
)

_STATUSCODE = _descriptor.EnumDescriptor(
  name='StatusCode',
  full_name='faiss.StatusCode',
  filename=None,
  file=DESCRIPTOR,
  create_key=_descriptor._internal_create_key,
  values=[
    _descriptor.EnumValueDescriptor(
      name='ENGINE_ERROR', index=0, number=0,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='SUCCESS', index=1, number=1,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='BAD_IMAGE', index=2, number=2,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='BAD_VECTOR', index=3, number=3,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='IMAGE_REGISTERED_FOR_ANOTHER_USER', index=4, number=4,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='FOUND', index=5, number=5,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='NOT_FOUND', index=6, number=6,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=553,
  serialized_end=692,
)
_sym_db.RegisterEnumDescriptor(_STATUSCODE)

StatusCode = enum_type_wrapper.EnumTypeWrapper(_STATUSCODE)
ENGINE_ERROR = 0
SUCCESS = 1
BAD_IMAGE = 2
BAD_VECTOR = 3
IMAGE_REGISTERED_FOR_ANOTHER_USER = 4
FOUND = 5
NOT_FOUND = 6



_VECTOR = _descriptor.Descriptor(
  name='Vector',
  full_name='faiss.Vector',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='vector_id', full_name='faiss.Vector.vector_id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='data', full_name='faiss.Vector.data', index=1,
      number=2, type=1, cpp_type=5, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=36,
  serialized_end=77,
)


_VECTORENROLLMENTREQUEST = _descriptor.Descriptor(
  name='VectorEnrollmentRequest',
  full_name='faiss.VectorEnrollmentRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='user_uid', full_name='faiss.VectorEnrollmentRequest.user_uid', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='vectors', full_name='faiss.VectorEnrollmentRequest.vectors', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=79,
  serialized_end=154,
)


_VECTORDELETIONREQUEST = _descriptor.Descriptor(
  name='VectorDeletionRequest',
  full_name='faiss.VectorDeletionRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='image_ids', full_name='faiss.VectorDeletionRequest.image_ids', index=0,
      number=1, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='user_uid', full_name='faiss.VectorDeletionRequest.user_uid', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=156,
  serialized_end=216,
)


_VECTORENROLLMENTRESPONSE = _descriptor.Descriptor(
  name='VectorEnrollmentResponse',
  full_name='faiss.VectorEnrollmentResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='code', full_name='faiss.VectorEnrollmentResponse.code', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='message', full_name='faiss.VectorEnrollmentResponse.message', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=218,
  serialized_end=294,
)


_VECTORSEARCHREQUEST = _descriptor.Descriptor(
  name='VectorSearchRequest',
  full_name='faiss.VectorSearchRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='vectors', full_name='faiss.VectorSearchRequest.vectors', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=296,
  serialized_end=349,
)


_VECTORSEARCHRESPONSE = _descriptor.Descriptor(
  name='VectorSearchResponse',
  full_name='faiss.VectorSearchResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='code', full_name='faiss.VectorSearchResponse.code', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='message', full_name='faiss.VectorSearchResponse.message', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='user_uid', full_name='faiss.VectorSearchResponse.user_uid', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='image_id', full_name='faiss.VectorSearchResponse.image_id', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='score', full_name='faiss.VectorSearchResponse.score', index=4,
      number=5, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=351,
  serialized_end=474,
)


_VECTORDELETIONRESPONSE = _descriptor.Descriptor(
  name='VectorDeletionResponse',
  full_name='faiss.VectorDeletionResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='code', full_name='faiss.VectorDeletionResponse.code', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='message', full_name='faiss.VectorDeletionResponse.message', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=476,
  serialized_end=550,
)

_VECTORENROLLMENTREQUEST.fields_by_name['vectors'].message_type = _VECTOR
_VECTORENROLLMENTRESPONSE.fields_by_name['code'].enum_type = _STATUSCODE
_VECTORSEARCHREQUEST.fields_by_name['vectors'].message_type = _VECTOR
_VECTORSEARCHRESPONSE.fields_by_name['code'].enum_type = _STATUSCODE
_VECTORDELETIONRESPONSE.fields_by_name['code'].enum_type = _STATUSCODE
DESCRIPTOR.message_types_by_name['Vector'] = _VECTOR
DESCRIPTOR.message_types_by_name['VectorEnrollmentRequest'] = _VECTORENROLLMENTREQUEST
DESCRIPTOR.message_types_by_name['VectorDeletionRequest'] = _VECTORDELETIONREQUEST
DESCRIPTOR.message_types_by_name['VectorEnrollmentResponse'] = _VECTORENROLLMENTRESPONSE
DESCRIPTOR.message_types_by_name['VectorSearchRequest'] = _VECTORSEARCHREQUEST
DESCRIPTOR.message_types_by_name['VectorSearchResponse'] = _VECTORSEARCHRESPONSE
DESCRIPTOR.message_types_by_name['VectorDeletionResponse'] = _VECTORDELETIONRESPONSE
DESCRIPTOR.enum_types_by_name['StatusCode'] = _STATUSCODE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Vector = _reflection.GeneratedProtocolMessageType('Vector', (_message.Message,), {
  'DESCRIPTOR' : _VECTOR,
  '__module__' : 'app.main.grpc.faiss_pb2'
  # @@protoc_insertion_point(class_scope:faiss.Vector)
  })
_sym_db.RegisterMessage(Vector)

VectorEnrollmentRequest = _reflection.GeneratedProtocolMessageType('VectorEnrollmentRequest', (_message.Message,), {
  'DESCRIPTOR' : _VECTORENROLLMENTREQUEST,
  '__module__' : 'app.main.grpc.faiss_pb2'
  # @@protoc_insertion_point(class_scope:faiss.VectorEnrollmentRequest)
  })
_sym_db.RegisterMessage(VectorEnrollmentRequest)

VectorDeletionRequest = _reflection.GeneratedProtocolMessageType('VectorDeletionRequest', (_message.Message,), {
  'DESCRIPTOR' : _VECTORDELETIONREQUEST,
  '__module__' : 'app.main.grpc.faiss_pb2'
  # @@protoc_insertion_point(class_scope:faiss.VectorDeletionRequest)
  })
_sym_db.RegisterMessage(VectorDeletionRequest)

VectorEnrollmentResponse = _reflection.GeneratedProtocolMessageType('VectorEnrollmentResponse', (_message.Message,), {
  'DESCRIPTOR' : _VECTORENROLLMENTRESPONSE,
  '__module__' : 'app.main.grpc.faiss_pb2'
  # @@protoc_insertion_point(class_scope:faiss.VectorEnrollmentResponse)
  })
_sym_db.RegisterMessage(VectorEnrollmentResponse)

VectorSearchRequest = _reflection.GeneratedProtocolMessageType('VectorSearchRequest', (_message.Message,), {
  'DESCRIPTOR' : _VECTORSEARCHREQUEST,
  '__module__' : 'app.main.grpc.faiss_pb2'
  # @@protoc_insertion_point(class_scope:faiss.VectorSearchRequest)
  })
_sym_db.RegisterMessage(VectorSearchRequest)

VectorSearchResponse = _reflection.GeneratedProtocolMessageType('VectorSearchResponse', (_message.Message,), {
  'DESCRIPTOR' : _VECTORSEARCHRESPONSE,
  '__module__' : 'app.main.grpc.faiss_pb2'
  # @@protoc_insertion_point(class_scope:faiss.VectorSearchResponse)
  })
_sym_db.RegisterMessage(VectorSearchResponse)

VectorDeletionResponse = _reflection.GeneratedProtocolMessageType('VectorDeletionResponse', (_message.Message,), {
  'DESCRIPTOR' : _VECTORDELETIONRESPONSE,
  '__module__' : 'app.main.grpc.faiss_pb2'
  # @@protoc_insertion_point(class_scope:faiss.VectorDeletionResponse)
  })
_sym_db.RegisterMessage(VectorDeletionResponse)


DESCRIPTOR._options = None

_VECTORAPI = _descriptor.ServiceDescriptor(
  name='VectorAPI',
  full_name='faiss.VectorAPI',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=695,
  serialized_end=919,
  methods=[
  _descriptor.MethodDescriptor(
    name='Enroll',
    full_name='faiss.VectorAPI.Enroll',
    index=0,
    containing_service=None,
    input_type=_VECTORENROLLMENTREQUEST,
    output_type=_VECTORENROLLMENTRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Search',
    full_name='faiss.VectorAPI.Search',
    index=1,
    containing_service=None,
    input_type=_VECTORSEARCHREQUEST,
    output_type=_VECTORSEARCHRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='Delete',
    full_name='faiss.VectorAPI.Delete',
    index=2,
    containing_service=None,
    input_type=_VECTORDELETIONREQUEST,
    output_type=_VECTORDELETIONRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_VECTORAPI)

DESCRIPTOR.services_by_name['VectorAPI'] = _VECTORAPI

# @@protoc_insertion_point(module_scope)
