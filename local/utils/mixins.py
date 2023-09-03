class UseSingleton:
    _instances_count = 0

    def __new__(cls):
        cls._instances_count += 1
        if not hasattr(cls, "instance"):
            cls.instance = super().__new__(cls)
        return cls.instance
