import 'package:property_change_notifier/property_change_notifier.dart';

enum CompareSchejControllerProperties {
  userIds,
  includeSelf,
  activeUserId,
}

class CompareSchejController extends PropertyChangeNotifier {
  final Set<String>? initialUserIds;
  final bool initialIncludeSelf;
  final String? initialActiveUserId;

  final Set<String> _userIds = <String>{};
  String? _activeUserId;
  late bool _includeSelf;

  CompareSchejController({
    this.initialUserIds,
    this.initialIncludeSelf = true,
    this.initialActiveUserId,
  }) {
    if (initialUserIds != null) _userIds.addAll(initialUserIds!);
    _includeSelf = initialIncludeSelf;
    _activeUserId = initialActiveUserId;
  }

  Set<String> get userIds => _userIds;

  bool get includeSelf => _includeSelf;
  set includeSelf(bool value) {
    _includeSelf = value;
    notifyListeners(CompareSchejControllerProperties.includeSelf);
  }

  String? get activeUserId => _activeUserId;
  set activeUserId(String? value) {
    _activeUserId = value;
    notifyListeners(CompareSchejControllerProperties.activeUserId);
  }

  void addUserId(String id) {
    _userIds.add(id);
    notifyListeners(CompareSchejControllerProperties.userIds);
  }

  void removeUserId(String id) {
    _userIds.remove(id);
    notifyListeners(CompareSchejControllerProperties.userIds);
  }
}
