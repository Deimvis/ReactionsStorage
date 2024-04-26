import logging
import shutil
from pathlib import Path
from typing import Callable


def find_impl_dir(dir_path: Path, impl_name: str) -> Path | None:
    for child in dir_path.iterdir():
        if child.is_dir() and child.name == f'.impl_{impl_name}':
            return child
    return None


def is_hidden(p: Path):
    return p.name.startswith('.')


def remove(p: Path):
    if p.is_dir():
        shutil.rmtree(p)
    else:
        p.unlink()


def copy(src: Path, dst_dir: Path):
    if src.is_dir():
        shutil.copytree(src, dst_dir / src.name)
    else:
        shutil.copy(src, dst_dir)
        

def clean_directory(dir_path: Path):
    """ ignores hidden files """
    for child in dir_path.iterdir():
        if not is_hidden(child):
            logging.info(f'remove: {child.relative_to(Path.cwd())}')
            remove(child)


def copy_directory(src_dir: Path, dst_dir: Path):
    """ ignores hidden files """
    for child in src_dir.iterdir():
        if not is_hidden(child):
            logging.info(f'copy: {child.relative_to(Path.cwd())} -> {dst_dir.relative_to(Path.cwd())}/')
            copy(child, dst_dir)
            

def traverse(dir_path: Path, target_impl_name: str, action: Callable[[Path, Path], None]):
    impl_dir = find_impl_dir(dir_path, target_impl_name)
    if impl_dir is not None:
        logging.info(f'Found impl in directory: {impl_dir.relative_to(Path.cwd())}')
        action(dir_path, impl_dir)
    else:
        for child in dir_path.iterdir():
            if child.is_dir():
                traverse(child, target_impl_name, action)
