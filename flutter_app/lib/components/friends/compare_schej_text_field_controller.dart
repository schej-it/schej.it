import 'package:flutter/material.dart';

class CompareSchejTextFieldController extends ChangeNotifier {
  final Set<String>? initialUserIds;
  final bool initialIncludeSelf;

  final Set<String> _userIds = <String>{};
  late bool _includeSelf;

  CompareSchejTextFieldController({
    this.initialUserIds,
    this.initialIncludeSelf = true,
  }) {
    if (initialUserIds != null) _userIds.addAll(initialUserIds!);
    _includeSelf = initialIncludeSelf;
  }

  Set<String> get userIds => _userIds;

  bool get includeSelf => _includeSelf;
  set includeSelf(bool value) {
    _includeSelf = value;
    notifyListeners();
  }

  void addUserId(String id) {
    _userIds.add(id);
    notifyListeners();
  }

  void removeUserId(String id) {
    _userIds.remove(id);
    notifyListeners();
  }
}
