import datetime

from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class JobState(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    JOB_STATE_UNSPECIFIED: _ClassVar[JobState]
    JOB_STATE_PENDING: _ClassVar[JobState]
    JOB_STATE_RUNNING: _ClassVar[JobState]
    JOB_STATE_SUCCEEDED: _ClassVar[JobState]
    JOB_STATE_FAILED: _ClassVar[JobState]
    JOB_STATE_CANCELLED: _ClassVar[JobState]
JOB_STATE_UNSPECIFIED: JobState
JOB_STATE_PENDING: JobState
JOB_STATE_RUNNING: JobState
JOB_STATE_SUCCEEDED: JobState
JOB_STATE_FAILED: JobState
JOB_STATE_CANCELLED: JobState

class Bus(_message.Message):
    __slots__ = ("id", "injection_mw")
    ID_FIELD_NUMBER: _ClassVar[int]
    INJECTION_MW_FIELD_NUMBER: _ClassVar[int]
    id: int
    injection_mw: float
    def __init__(self, id: _Optional[int] = ..., injection_mw: _Optional[float] = ...) -> None: ...

class Line(_message.Message):
    __slots__ = ("from_bus", "to_bus", "reactance_pu")
    FROM_BUS_FIELD_NUMBER: _ClassVar[int]
    TO_BUS_FIELD_NUMBER: _ClassVar[int]
    REACTANCE_PU_FIELD_NUMBER: _ClassVar[int]
    from_bus: int
    to_bus: int
    reactance_pu: float
    def __init__(self, from_bus: _Optional[int] = ..., to_bus: _Optional[int] = ..., reactance_pu: _Optional[float] = ...) -> None: ...

class DcPowerFlowSpec(_message.Message):
    __slots__ = ("buses", "lines", "slack_bus_id", "base_mva")
    BUSES_FIELD_NUMBER: _ClassVar[int]
    LINES_FIELD_NUMBER: _ClassVar[int]
    SLACK_BUS_ID_FIELD_NUMBER: _ClassVar[int]
    BASE_MVA_FIELD_NUMBER: _ClassVar[int]
    buses: _containers.RepeatedCompositeFieldContainer[Bus]
    lines: _containers.RepeatedCompositeFieldContainer[Line]
    slack_bus_id: int
    base_mva: float
    def __init__(self, buses: _Optional[_Iterable[_Union[Bus, _Mapping]]] = ..., lines: _Optional[_Iterable[_Union[Line, _Mapping]]] = ..., slack_bus_id: _Optional[int] = ..., base_mva: _Optional[float] = ...) -> None: ...

class LineFlow(_message.Message):
    __slots__ = ("from_bus", "to_bus", "flow_mw")
    FROM_BUS_FIELD_NUMBER: _ClassVar[int]
    TO_BUS_FIELD_NUMBER: _ClassVar[int]
    FLOW_MW_FIELD_NUMBER: _ClassVar[int]
    from_bus: int
    to_bus: int
    flow_mw: float
    def __init__(self, from_bus: _Optional[int] = ..., to_bus: _Optional[int] = ..., flow_mw: _Optional[float] = ...) -> None: ...

class DcPowerFlowResult(_message.Message):
    __slots__ = ("bus_angles_rad", "line_flows", "slack_injection_mw")
    class BusAnglesRadEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: int
        value: float
        def __init__(self, key: _Optional[int] = ..., value: _Optional[float] = ...) -> None: ...
    BUS_ANGLES_RAD_FIELD_NUMBER: _ClassVar[int]
    LINE_FLOWS_FIELD_NUMBER: _ClassVar[int]
    SLACK_INJECTION_MW_FIELD_NUMBER: _ClassVar[int]
    bus_angles_rad: _containers.ScalarMap[int, float]
    line_flows: _containers.RepeatedCompositeFieldContainer[LineFlow]
    slack_injection_mw: float
    def __init__(self, bus_angles_rad: _Optional[_Mapping[int, float]] = ..., line_flows: _Optional[_Iterable[_Union[LineFlow, _Mapping]]] = ..., slack_injection_mw: _Optional[float] = ...) -> None: ...

class JobSpec(_message.Message):
    __slots__ = ("dc_power_flow", "simulated_delay_ms")
    DC_POWER_FLOW_FIELD_NUMBER: _ClassVar[int]
    SIMULATED_DELAY_MS_FIELD_NUMBER: _ClassVar[int]
    dc_power_flow: DcPowerFlowSpec
    simulated_delay_ms: int
    def __init__(self, dc_power_flow: _Optional[_Union[DcPowerFlowSpec, _Mapping]] = ..., simulated_delay_ms: _Optional[int] = ...) -> None: ...

class JobResult(_message.Message):
    __slots__ = ("dc_power_flow",)
    DC_POWER_FLOW_FIELD_NUMBER: _ClassVar[int]
    dc_power_flow: DcPowerFlowResult
    def __init__(self, dc_power_flow: _Optional[_Union[DcPowerFlowResult, _Mapping]] = ...) -> None: ...

class Job(_message.Message):
    __slots__ = ("id", "state", "spec", "result", "error_message", "created_at", "started_at", "finished_at")
    ID_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    SPEC_FIELD_NUMBER: _ClassVar[int]
    RESULT_FIELD_NUMBER: _ClassVar[int]
    ERROR_MESSAGE_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    STARTED_AT_FIELD_NUMBER: _ClassVar[int]
    FINISHED_AT_FIELD_NUMBER: _ClassVar[int]
    id: str
    state: JobState
    spec: JobSpec
    result: JobResult
    error_message: str
    created_at: _timestamp_pb2.Timestamp
    started_at: _timestamp_pb2.Timestamp
    finished_at: _timestamp_pb2.Timestamp
    def __init__(self, id: _Optional[str] = ..., state: _Optional[_Union[JobState, str]] = ..., spec: _Optional[_Union[JobSpec, _Mapping]] = ..., result: _Optional[_Union[JobResult, _Mapping]] = ..., error_message: _Optional[str] = ..., created_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., started_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., finished_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ...) -> None: ...

class SubmitJobRequest(_message.Message):
    __slots__ = ("spec",)
    SPEC_FIELD_NUMBER: _ClassVar[int]
    spec: JobSpec
    def __init__(self, spec: _Optional[_Union[JobSpec, _Mapping]] = ...) -> None: ...

class SubmitJobResponse(_message.Message):
    __slots__ = ("job_id",)
    JOB_ID_FIELD_NUMBER: _ClassVar[int]
    job_id: str
    def __init__(self, job_id: _Optional[str] = ...) -> None: ...

class GetJobRequest(_message.Message):
    __slots__ = ("job_id",)
    JOB_ID_FIELD_NUMBER: _ClassVar[int]
    job_id: str
    def __init__(self, job_id: _Optional[str] = ...) -> None: ...

class GetJobResponse(_message.Message):
    __slots__ = ("job",)
    JOB_FIELD_NUMBER: _ClassVar[int]
    job: Job
    def __init__(self, job: _Optional[_Union[Job, _Mapping]] = ...) -> None: ...

class WatchJobRequest(_message.Message):
    __slots__ = ("job_id",)
    JOB_ID_FIELD_NUMBER: _ClassVar[int]
    job_id: str
    def __init__(self, job_id: _Optional[str] = ...) -> None: ...

class WatchJobResponse(_message.Message):
    __slots__ = ("job",)
    JOB_FIELD_NUMBER: _ClassVar[int]
    job: Job
    def __init__(self, job: _Optional[_Union[Job, _Mapping]] = ...) -> None: ...

class CancelJobRequest(_message.Message):
    __slots__ = ("job_id",)
    JOB_ID_FIELD_NUMBER: _ClassVar[int]
    job_id: str
    def __init__(self, job_id: _Optional[str] = ...) -> None: ...

class CancelJobResponse(_message.Message):
    __slots__ = ("job",)
    JOB_FIELD_NUMBER: _ClassVar[int]
    job: Job
    def __init__(self, job: _Optional[_Union[Job, _Mapping]] = ...) -> None: ...
