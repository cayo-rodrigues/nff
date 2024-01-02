import os


class FileManager:
    @classmethod
    def file_exists(cls, path: str) -> bool:
        return os.path.isfile(path)

    @classmethod
    def dir_exists(cls, path: str) -> bool:
        return os.path.isdir(path)

    @classmethod
    def get_or_create_dir(cls, path: str) -> str:
        if not cls.dir_exists(path):
            os.mkdir(path)
        return path

    @classmethod
    def list_file_names(cls, dir_path: str) -> list[str]:
        return os.listdir(cls.get_or_create_dir(dir_path))

    @classmethod
    def count_files(cls, dir_path: str) -> int:
        return len(cls.list_file_names(dir_path))

    @classmethod
    def rename_file(cls, old_name: str, new_name: str) -> None:
        if cls.file_exists(old_name):
            os.rename(src=old_name, dst=new_name)

    @classmethod
    def get_latest_file_name(cls, dir_path: str) -> str:
        file_names = cls.list_file_names(dir_path)
        if not file_names:
            return ""

        return max(
            [dir_path + name for name in file_names],
            key=os.path.getctime,
        )

    @classmethod
    def get_file_name_from_path(cls, path: str) -> str:
        return os.path.basename(path)

    @classmethod
    def erase_file(cls, path: str) -> None:
        if cls.file_exists(path):
            os.remove(path)
